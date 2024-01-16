package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/gosip/sdp"
	"github.com/i-Things/things/src/vidsvr/gosip/sip"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// 从请求中解析出设备信息
func parserDevicesFromReqeust(req *sip.Request) (GbSipDevice, bool) {
	u := GbSipDevice{}
	header, ok := req.From()
	if !ok {
		logrus.Warningln("not found from header from request", req.String())
		return u, false
	}
	if header.Address == nil {
		logrus.Warningln("not found from user from request", req.String())
		return u, false
	}
	if header.Address.User() == nil {
		logrus.Warningln("not found from user from request", req.String())
		return u, false
	}
	u.DeviceID = header.Address.User().String()
	u.Region = header.Address.Host()
	via, ok := req.ViaHop()
	if !ok {
		logrus.Info("not found ViaHop from request", req.String())
		return u, false
	}
	u.Host = via.Host
	u.Port = via.Port.String()
	report, ok := via.Params.Get("rport")
	if ok && report != nil {
		u.Rport = report.String()
	}
	raddr, ok := via.Params.Get("received")
	if ok && raddr != nil {
		u.RAddr = raddr.String()
	}

	u.TransPort = via.Transport
	u.URIStr = header.Address.String()
	//byteTmp, _ := json.Marshal()
	fmt.Println("[airgens-sip] --------:", sip.NewAddressFromFromHeader(header))
	u.addr = sip.NewAddressFromFromHeader(header)
	u.Source = req.Source().String()
	u.source = req.Source()
	return u, true
}

// 获取设备信息（注册设备）
func sipDeviceInfo(to GbSipDevice) {

	hb := sip.NewHeaderBuilder().SetTo(to.addr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE)
	req := sip.NewRequest("", sip.MESSAGE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetDeviceInfoXML(to.DeviceID))
	req.SetDestination(to.source)
	tx, err := SipSrv.Srv.Request(req)
	if err != nil {
		logrus.Warnln("sipDeviceInfo  error,", err)
		return
	}
	_, err = sipResponse(tx)
	if err != nil {
		logrus.Warnln("sipDeviceInfo  response error,", err)
		return
	}
}

func sipResponse(tx *sip.Transaction) (*sip.Response, error) {
	response := tx.GetResponse()
	if response == nil {
		return nil, utils.NewError(nil, "response timeout", "tx key:", tx.Key())
	}
	if response.StatusCode() != http.StatusOK {
		return response, utils.NewError(nil, "response fail", response.StatusCode(), response.Reason(), "tx key:", tx.Key())
	}
	return response, nil
}

func sipMessageCatalog(body []byte) error {
	message := &MessageDeviceListResponse{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	channelsRepo := db.NewVidmgrChannelsRepo(Ctx)
	fmt.Println("[airgens-sip] sipMessageCatalog-message:", message)
	if message.SumNum > 0 {
		fmt.Println("[airgens-sip] message:", message)
		for _, d := range message.Item {
			filter := db.VidmgrChannelsFilter{
				DeviceIDs: []string{message.DeviceID},
				ChannelID: d.ChannelID,
			}
			tmpStr, _ := json.Marshal(d)

			fmt.Println("[airgens-sip] filter_SIP devices:", filter)
			fmt.Println("[airgens-sip] filter_SIP Item-d:", d)
			fmt.Println("[airgens-sip] filter_SIP d.ChannelID:", d.ChannelID)
			fmt.Println("[airgens-sip] filter_SIP json.tmpStr:", string(tmpStr))
			channel, err := channelsRepo.FindOneByFilter(Ctx, filter)
			if err != nil {
				errosType := &types.IndexApiResp{}
				json.Unmarshal([]byte(err.Error()), errosType)
				//not found  and set default
				if errosType.Code == 100009 {
					channel = &db.VidmgrChannels{}
					channel.DeviceID = message.DeviceID
					channel.ChannelID = d.ChannelID
					channel.Active = time.Now().Unix()
					channel.URIStr = fmt.Sprintf("sip:%s@%s", d.ChannelID, SipInfo.Region)
					channel.Status = transDeviceStatus(d.Status)
					channel.Name = d.Name
					channel.Manufacturer = d.Manufacturer
					channel.Model = d.Model
					channel.Owner = d.Owner
					channel.CivilCode = d.CivilCode
					// Address ip地址
					channel.Address = d.Address
					channel.Parental = d.Parental
					channel.SafetyWay = d.SafetyWay
					channel.RegisterWay = d.RegisterWay
					channel.Secrecy = d.Secrecy
					channel.LastLogin = time.Now()
					channelsRepo.Insert(Ctx, channel)
				}
			} else { //正常查询到设备
				channel.Active = time.Now().Unix()
				channel.URIStr = fmt.Sprintf("sip:%s@%s", d.ChannelID, SipInfo.Region)
				channel.Status = transDeviceStatus(d.Status)
				channel.Name = d.Name
				channel.Manufacturer = d.Manufacturer
				channel.Model = d.Model
				channel.Owner = d.Owner
				channel.CivilCode = d.CivilCode
				// Address ip地址
				channel.Address = d.Address
				channel.Parental = d.Parental
				channel.SafetyWay = d.SafetyWay
				channel.RegisterWay = d.RegisterWay
				channel.Secrecy = d.Secrecy
				channelsRepo.Update(Ctx, channel)
			}
			go notify(notifyChannelsActive(*channel))
		}
	}
	return nil
}
func sipMessageKeepalive(u GbSipDevice, body []byte) error {
	message := &MessageNotify{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	deviceRepo := db.NewVidmgrDevicesRepo(Ctx)
	filter := db.VidmgrDevicesFilter{
		DeviceID: u.DeviceID,
	}
	device, err := deviceRepo.FindOneByFilter(Ctx, filter)
	if err != nil {
		logrus.Warnln("Device Keepalive not found ", u.DeviceID, err)
		return err
	}
	if device != nil {
		device.Host = u.Host
		device.Port = u.Port
		device.Rport = u.Rport
		device.RAddr = u.RAddr
		device.Source = u.Source
		device.URIStr = u.URIStr
		if message.Status == "OK" {
			device.LastLogin = time.Now()
		}
		err = deviceRepo.Update(Ctx, device)
	}
	go notify(notifyDevicesAcitve(device.DeviceID, message.Status))

	return err
}

func sipMessageRecordInfo(u GbSipDevice, body []byte) error {
	message := &MessageRecordInfoResponse{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	recordKey := fmt.Sprintf("%s%d", message.DeviceID, message.SN)
	if list, ok := _recordList.Load(recordKey); ok {
		info := list.(recordList)
		info.l.Lock()
		defer info.l.Unlock()
		info.num += len(message.Item)
		var sint, eint int64
		for _, item := range message.Item {
			s, _ := time.ParseInLocation("2006-01-02T15:04:05", item.StartTime, time.Local)
			e, _ := time.ParseInLocation("2006-01-02T15:04:05", item.EndTime, time.Local)
			sint = s.Unix()
			eint = e.Unix()
			if sint < info.s {
				sint = info.s
			}
			if eint > info.e {
				eint = info.e
			}
			info.data = append(info.data, []int64{sint, eint})
		}
		if info.num == message.SumNum {
			// 获取到完整数据
			info.resp <- transRecordList(info.data)
		}
		_recordList.Store(recordKey, info)
		return nil
	}
	return errors.MediaRecordNotFound.AddDetail("recordlist devices not found")
}

func sipMessageDeviceInfo(body []byte) error {
	message := &MessageDeviceInfoResponse{}
	if err := utils.XMLDecode([]byte(body), message); err != nil {
		logrus.Errorln("sipMessageDeviceInfo Unmarshal xml err:", err, "body:", body)
		return err
	}
	deviceRepo := db.NewVidmgrDevicesRepo(Ctx)
	dev, err := deviceRepo.FindOneByFilter(Ctx, db.VidmgrDevicesFilter{
		DeviceID: message.DeviceID,
	})
	if err != nil {
		fmt.Println("_____sipMessageDeviceInfo error___ not found data")
		return nil
	}
	dev.Model = message.Model
	dev.DeviceType = message.DeviceType
	dev.Firmware = message.Firmware
	dev.Manufacturer = message.Manufacturer
	dev.LastLogin = time.Now()
	deviceRepo.Update(Ctx, dev)
	return nil
}

// sipCatalog 获取注册设备包含的列表
func sipCatalog(to GbSipDevice) {
	hb := sip.NewHeaderBuilder().SetTo(to.addr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE)
	req := sip.NewRequest("", sip.MESSAGE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetCatalogXML(to.DeviceID))
	req.SetDestination(to.source)
	tx, err := SipSrv.Srv.Request(req)
	if err != nil {
		logrus.Warnln("sipCatalog  error,", err)
		return
	}
	_, err = sipResponse(tx)
	if err != nil {
		logrus.Warnln("sipCatalog  response error,", err)
		return
	}
}

func transDeviceStatus(status string) string {
	if v, ok := deviceStatusMap[status]; ok {
		return v
	}
	return status
}

// 将返回的多组数据合并，时间连续的进行合并，最后按照天返回数据，返回为某天内时间段列表
func transRecordList(data [][]int64) Records {
	if len(data) == 0 {
		return Records{}
	}
	res := Records{}
	list := map[string][]RecordInfo{}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})
	newData := [][]int64{}
	var newDataIE = []int64{}

	for x, d := range data {
		if x == 0 {
			newDataIE = d
			continue
		}
		if d[0] == newDataIE[1] {
			newDataIE[1] = d[1]
		} else {
			newData = append(newData, newDataIE)
			newDataIE = d
		}
	}
	newData = append(newData, newDataIE)
	var cs, ce time.Time
	dates := []string{}
	for _, d := range newData {
		s := time.Unix(d[0], 0)
		e := time.Unix(d[1], 0)
		cs, _ = time.ParseInLocation("20060102", s.Format("20060102"), time.Local)
		for {
			ce = cs.Add(24 * time.Hour)
			if e.Unix() >= ce.Unix() {
				// 当前时段跨天
				if v, ok := list[cs.Format("2006-01-02")]; ok {
					list[cs.Format("2006-01-02")] = append(v, RecordInfo{
						Start: utils.MaxTime(s.Unix(), cs.Unix()),
						End:   ce.Unix() - 1,
					})
				} else {
					list[cs.Format("2006-01-02")] = []RecordInfo{
						{
							Start: utils.MaxTime(s.Unix(), cs.Unix()),
							End:   ce.Unix() - 1,
						},
					}
					dates = append(dates, cs.Format("2006-01-02"))
					res.DayTotal++
				}
				res.TimeNum++
				cs = ce
			} else {
				if v, ok := list[cs.Format("2006-01-02")]; ok {
					list[cs.Format("2006-01-02")] = append(v, RecordInfo{
						Start: utils.MaxTime(s.Unix(), cs.Unix()),
						End:   e.Unix(),
					})
				} else {
					list[cs.Format("2006-01-02")] = []RecordInfo{
						{
							Start: utils.MaxTime(s.Unix(), cs.Unix()),
							End:   e.Unix(),
						},
					}
					dates = append(dates, cs.Format("2006-01-02"))
					res.DayTotal++
				}
				res.TimeNum++
				break
			}
		}
	}
	resData := []RecordDate{}
	for _, date := range dates {
		resData = append(resData, RecordDate{
			Date:  date,
			Items: list[date],
		})

	}
	res.Data = resData
	return res
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func sipPlayPush(data *Stream) (*Stream, error) {
	var (
		s sdp.Session
		b []byte
	)
	name := "Play"
	protocal := "TCP/RTP/AVP"
	if data.Type == 1 { // 0  直播 1 历史
		name = "Playback"
		protocal = "RTP/RTCP"
	}
	video := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "video",
			Port:     data.MediaPort,
			Formats:  []string{"96", "98", "97"},
			Protocol: protocal,
		},
	}
	video.AddAttribute("recvonly")
	if data.Type == 0 {
		video.AddAttribute("setup", "passive")
		video.AddAttribute("connection", "new")
	}
	video.AddAttribute("rtpmap", "96", "PS/90000")
	video.AddAttribute("rtpmap", "98", "H264/90000")
	video.AddAttribute("rtpmap", "97", "MPEG4/90000")
	// defining message
	msg := &sdp.Message{
		Origin: sdp.Origin{
			Username: data.DeviceID, // 媒体服务器id
			Address:  data.MediaIP,
		},
		Name: name,
		Connection: sdp.ConnectionData{
			IP:  net.ParseIP(data.MediaIP),
			TTL: 0,
		},
		Timing: []sdp.Timing{
			{
				Start: time.Time{},
				End:   time.Time{},
			},
		},
		Medias: []sdp.Media{video},
		SSRC:   getSSRC(data.Type),
	}
	if data.Type == 1 {
		msg.URI = fmt.Sprintf("%s:0", data.ChannelID)
	}
	// appending message to session
	s = msg.Append(s)
	// appending session to byte buffer
	b = s.AppendTo(b)
	uri, _ := sip.ParseURI(data.ChnURIStr)
	ChnAddr := &sip.Address{URI: uri}
	_serverDevices.addr.Params.Add("tag", sip.String{Str: utils.RandString(20)})
	hb := sip.NewHeaderBuilder().SetTo(ChnAddr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeSDP).SetMethod(sip.INVITE).SetContact(_serverDevices.addr)
	req := sip.NewRequest("", sip.INVITE, ChnAddr.URI, sip.DefaultSipVersion, hb.Build(), b)
	Source, _ := net.ResolveUDPAddr("udp", data.DevSource)
	req.SetDestination(Source)
	req.AppendHeader(&sip.GenericHeader{HeaderName: "Subject", Contents: fmt.Sprintf("%s:%s,%s:%s", data.ChannelID, data.Stream, _serverDevices.DeviceID, data.Stream)})
	req.SetRecipient(ChnAddr.URI)
	tx, err := SipSrv.Srv.Request(req)
	if err != nil {
		logrus.Warningln("sipPlayPush fail.id:", data.DeviceID, data.ChannelID, "err:", err)
		return data, err
	}
	response, err := sipResponse(tx)
	if err != nil {
		logrus.Warningln("sipPlayPush response fail.id:", data.DeviceID, data.ChannelID, "err:", err)
		return data, err
	}
	tx.Request(sip.NewRequestFromResponse(sip.ACK, response))
	return data, err
}

// 这个地方需要再研究一下
func getSSRC(t int) string {
	return fmt.Sprintf("%d%s%04d", t, SipInfo.Region[3:8], _serverDevices.RandomStr.GetSnowflakeId()%10000)
}

// zlm接收到的ssrc为16进制。发起请求的ssrc为10进制
func ssrc2stream(ssrc string) string {
	if ssrc[0:1] == "0" {
		ssrc = ssrc[1:]
	}
	num, _ := strconv.Atoi(ssrc)
	return fmt.Sprintf("%08X", num)
}
