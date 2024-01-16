package media

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	sip "github.com/i-Things/things/src/vidsvr/gosip/sip"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func handlerRegister(req *sip.Request, tx *sip.Transaction) {
	// 判断是否存在授权字段
	if hdrs := req.GetHeaders("Authorization"); len(hdrs) > 0 {
		fmt.Println("--------------handlerRegister  Authorization-------------------")
		fromUser, ok := parserDevicesFromReqeust(req)
		if !ok {
			return
		}
		//查找该DvicesID
		deviceRepo := db.NewVidmgrDevicesRepo(Ctx)
		dev, err := deviceRepo.FindOneByFilter(Ctx, db.VidmgrDevicesFilter{
			DeviceID: fromUser.DeviceID,
		})
		if err != nil {
			//没有查到数据推出
			fmt.Println("____handlerRegister_errror____handlerRegister_errror:", err)
			//没有查到数据的话，直接进行录入数据库
			errosType := &types.IndexApiResp{}
			json.Unmarshal([]byte(err.Error()), errosType)
			//not found  and set default
			if errosType.Code == 100009 {
				dev = &db.VidmgrDevices{}
				fromUser.Name = Config.GbsipConf.DefaultDevName
				fromUser.PWD = Config.GbsipConf.DefaultDevPswd
				dev.PWD = fromUser.PWD
				dev.Name = fromUser.Name
				SipDeviceToDB(&fromUser, dev)
				dev.DeviceID = fromUser.DeviceID
				dev.LastLogin = time.Now()
				deviceRepo.Insert(Ctx, dev)
			} else {
				return
			}
		}
		//查到数据之后，对数据进行修改
		if dev != nil {
			if !dev.Regist {
				// 如果数据库里用户未激活，替换user数据
				fromUser.Name = dev.Name
				fromUser.PWD = dev.PWD
				//dev = &fromUser
				SipDeviceToDB(&fromUser, dev)
				dev.Regist = true
			}
			fmt.Println("[airgens-sip] fromUser:", fromUser)
			fmt.Println("[airgens-sip] dev:", dev)

			authenticateHeader := hdrs[0].(*sip.GenericHeader)
			auth := sip.AuthFromValue(authenticateHeader.Contents)
			auth.SetPassword(dev.PWD)
			auth.SetUsername(dev.DeviceID)
			auth.SetMethod(string(req.Method()))
			auth.SetURI(auth.Get("uri"))
			if auth.CalcResponse() == auth.Get("response") {
				// 验证成功
				// 记录活跃设备
				dev.LastLogin = time.Now()
				dev.Source = fromUser.source.String()
				//dev. = fromUser.addr
				deviceRepo.Update(Ctx, dev)
				tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
				// 注册成功后查询设备信息，获取制作厂商等信息
				go notify(notifyDevicesRegister(fromUser))
				go sipDeviceInfo(fromUser)
				return
			}
		}
	}
	resp := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
	resp.AppendHeader(&sip.GenericHeader{HeaderName: "WWW-Authenticate", Contents: fmt.Sprintf("Digest nonce=\"%s\", algorithm=MD5, realm=\"%s\",qop=\"auth\"", utils.RandString(32), SipInfo.Region)})
	tx.Respond(resp)
}

func handlerMessage(req *sip.Request, tx *sip.Transaction) {
	u, ok := parserDevicesFromReqeust(req)
	if !ok {
		// 未解析出来源用户返回错误
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
		return
	}
	// 判断是否存在body数据
	if len, have := req.ContentLength(); !have || len.Equals(0) {
		// 不存在就直接返回的成功
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	}
	body := req.Body()
	message := &MessageReceive{}

	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Warnln("Message Unmarshal xml err:", err, "body:", string(body))
		// 有些body xml发送过来的不带encoding ，而且格式不是utf8的，导致xml解析失败，此处使用gbk转utf8后再次尝试xml解析
		body, err = GbkToUtf8(body)
		if err != nil {
			logrus.Errorln("message gbk to utf8 err", err)
		}
		if err := utils.XMLDecode(body, message); err != nil {
			logrus.Errorln("Message Unmarshal xml after gbktoutf8 err:", err, "body:", string(body))
			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
			return
		}
	}
	switch message.CmdType {
	case "Catalog":
		// 设备列表
		fmt.Println("--------------handlerMessage  Catalog 通道登记-------------------")
		sipMessageCatalog(body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	case "Keepalive":
		// heardbeat
		fmt.Println("--------------handlerMessage  Keepalive-------------------")
		if err := sipMessageKeepalive(u, body); err == nil {
			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
			// 心跳后同步注册设备列表信息
			sipCatalog(u)
			return
		}
	case "RecordInfo":
		fmt.Println("--------------handlerMessage  RecordInfo-------------------")
		// 设备音视频文件列表
		sipMessageRecordInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	case "DeviceInfo":
		// 主设备信息
		fmt.Println("--------------handlerMessage  DeviceInfo-------------------")
		sipMessageDeviceInfo(body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	}
	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
}
