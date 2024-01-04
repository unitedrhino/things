package media

import (
	"fmt"
	"github.com/i-Things/things/shared/utils"
	sip "github.com/i-Things/things/src/vidsvr/gosip/sip"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/sirupsen/logrus"
	"net/http"
)

func handlerRegister(req *sip.Request, tx *sip.Transaction) {
	// 判断是否存在授权字段
	if hdrs := req.GetHeaders("Authorization"); len(hdrs) > 0 {
		fromUser, ok := parserDevicesFromReqeust(req)
		if !ok {
			return
		}
		//查找该DvicesID
		deviceRepo := db.NewVidmgrDevicesRepo(Ctx)
		user, err := deviceRepo.FindOneByFilter(Ctx, db.VidmgrDevicesFilter{
			DeviceIDs: []string{fromUser.DeviceID},
		})
		if err == nil {
			if !user.Regist {
				// 如果数据库里用户未激活，替换user数据
				fromUser.ID = user.ID
				fromUser.Name = user.Name
				fromUser.PWD = user.PWD
				user = &fromUser
			}
			user.Addr = fromUser.Addr
			authenticateHeader := hdrs[0].(*sip.GenericHeader)
			auth := sip.AuthFromValue(authenticateHeader.Contents)
			auth.SetPassword(user.PWD)
			auth.SetUsername(user.DeviceID)
			auth.SetMethod(string(req.Method()))
			auth.SetURI(auth.Get("uri"))

			if auth.CalcResponse() == auth.Get("response") {
				// 验证成功
				// 记录活跃设备
				user.Source = fromUser.Source
				user.Addr = fromUser.Addr
				//
				_activeDevices.Store(user.DeviceID, user)
				if !user.Regist {
					// 第一次激活，保存数据库
					user.Regist = true
					deviceRepo.Insert(Ctx, user)
					//db.DBClient.Save(&user)
					logrus.Infoln("new user regist,id:", user.DeviceID)
				}
				tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
				// 注册成功后查询设备信息，获取制作厂商等信息
				go notify(notifyDevicesRegister(*user))
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
		sipMessageCatalog(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	case "Keepalive":
		// heardbeat
		if err := sipMessageKeepalive(u, body); err == nil {
			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
			// 心跳后同步注册设备列表信息
			sipCatalog(u)
			return
		}
	case "RecordInfo":
		// 设备音视频文件列表
		sipMessageRecordInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	case "DeviceInfo":
		// 主设备信息
		sipMessageDeviceInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	}
	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
}
