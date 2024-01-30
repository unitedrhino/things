package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	sdp2 "github.com/i-Things/things/src/vidsip/gosip/sdp"
	sip2 "github.com/i-Things/things/src/vidsip/gosip/sip"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

var serverDevice db.SipDevices

// 从请求中解析出设备信息
func parserDevicesFromReqeust(req *sip2.Request) (db.SipDevices, bool) {
	u := db.SipDevices{}
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
	//fmt.Println("[airgens-sip] --------:", sip2.NewAddressFromFromHeader(header))
	u.Taddr = sip2.NewAddressFromFromHeader(header)
	u.Source = req.Source().String()
	u.Tsource = req.Source()
	return u, true
}

// 获取设备信息（注册设备）
func sipDeviceInfo(to db.SipDevices) {
	hb := sip2.NewHeaderBuilder().SetTo(to.Taddr).SetFrom(serverDevice.Taddr).AddVia(&sip2.ViaHop{
		Params: sip2.NewParams().Add("branch", sip2.String{Str: sip2.GenerateBranch()}),
	}).SetContentType(&sip2.ContentTypeXML).SetMethod(sip2.MESSAGE)
	req := sip2.NewRequest("", sip2.MESSAGE, to.Taddr.URI, sip2.DefaultSipVersion, hb.Build(), sip2.GetDeviceInfoXML(to.DeviceID))
	req.SetDestination(to.Tsource)
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

func sipResponse(tx *sip2.Transaction) (*sip2.Response, error) {
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
	channelsRepo := db.NewSipChannelsRepo(Ctx)
	fmt.Println("[airgens-sip] sipMessageCatalog-message:", message)
	if message.SumNum > 0 {
		for _, d := range message.Item {
			filter := db.SipChannelsFilter{
				DeviceIDs: []string{message.DeviceID},
				ChannelID: d.ChannelID,
			}
			tmpStr, _ := json.Marshal(d)

			fmt.Println("[airgens-sip] filter_SIP devices:", filter)
			fmt.Println("[airgens-sip] filter_SIP Item-d:", d)
			fmt.Println("[airgens-sip] filter_SIP d.ChannelID:", d.ChannelID)
			fmt.Println("[airgens-sip] filter_SIP json.tmpStr:", string(tmpStr))
			channel, err := channelsRepo.FindOneByFilter(Ctx, filter)
			if err == nil { //正常查询到设备
				channel.LastLogin = time.Now().Unix()
				channel.URIStr = fmt.Sprintf("sip:%s@%s", d.ChannelID, SipInfo.Region)
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
				channel.Status = transDeviceStatus(d.Status)
				channelsRepo.Update(Ctx, channel)
				go notify(notifyChannelsActive(*channel))
			} else {
				logrus.Infoln("deviceid not found,deviceid:", d.DeviceID, "pdid:", message.DeviceID, "err", err)
			}
		}
	}
	return nil
}
func sipMessageKeepalive(u db.SipDevices, body []byte) error {
	message := &MessageNotify{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	devRepo := db.NewSipDevicesRepo(Ctx)
	device, ok := activeDevices.Get(u.DeviceID)
	if !ok {

		filter := db.SipDevicesFilter{
			DeviceID: u.DeviceID,
		}
		dev, err := devRepo.FindOneByFilter(Ctx, filter)
		if err != nil {
			logrus.Warnln("DeviceID Keepalive not found ", u.DeviceID, err)
		}
		device = *dev
	}
	if message.Status == "OK" {
		device.LastLogin = time.Now().Unix()
		activeDevices.Store(u.DeviceID, u)
	} else {
		device.LastLogin = -1
		activeDevices.Delete(u.DeviceID)
	}
	go notify(notifyDevicesAcitve(device.DeviceID, message.Status))

	err := devRepo.Update(Ctx, &db.SipDevices{
		DeviceID:  u.DeviceID,
		Host:      u.Host,
		Port:      u.Port,
		Rport:     u.Rport,
		RAddr:     u.RAddr,
		Source:    u.Source,
		URIStr:    u.URIStr,
		LastLogin: device.LastLogin,
	})
	return err
}

func sipMessageRecordInfo(u db.SipDevices, body []byte) error {
	message := &MessageRecordInfoResponse{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	recordKey := fmt.Sprintf("%s%d", message.DeviceID, message.SN)
	if list, ok := recordList.Load(recordKey); ok {
		info := list.(RecordList)
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
		recordList.Store(recordKey, info)
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
	deviceRepo := db.NewSipDevicesRepo(Ctx)
	dev, err := deviceRepo.FindOneByFilter(Ctx, db.SipDevicesFilter{
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
	dev.LastLogin = time.Now().Unix()
	dev.Status = 1
	//if dev.VidmgrID == "" {
	//	dev.VidmgrID = SipInfo.MgrID
	//}
	deviceRepo.Update(Ctx, dev)
	return nil
}

// sipCatalog 获取注册设备包含的列表
func sipCatalog(to db.SipDevices) {
	hb := sip2.NewHeaderBuilder().SetTo(to.Taddr).SetFrom(serverDevice.Taddr).AddVia(&sip2.ViaHop{
		Params: sip2.NewParams().Add("branch", sip2.String{Str: sip2.GenerateBranch()}),
	}).SetContentType(&sip2.ContentTypeXML).SetMethod(sip2.MESSAGE)
	req := sip2.NewRequest("", sip2.MESSAGE, to.Taddr.URI, sip2.DefaultSipVersion, hb.Build(), sip2.GetCatalogXML(to.DeviceID))
	req.SetDestination(to.Tsource)
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

func Play(chnid string, replay int32, start int64, end int64) (any, error) {
	pm := &Stream{S: time.Time{}, E: time.Time{}, ChannelID: chnid}
	if replay == 1 {
		pm.Type = 1
		if start == 0 {
			return nil, errors.MediaSipPlayError.AddDetail("开始时间错误")
		}
		pm.S = time.Unix(start, 0)
		pm.E = time.Unix(end, 0)
		if start >= end {
			return nil, errors.MediaSipPlayError.AddDetail("开始时间>=结束时间")
		}
	} else {
		if succ, ok := StreamList.Succ.Load(chnid); ok {
			return succ, nil
		}
	}
	res, err := sipPlay(pm)
	if err != nil {
		return nil, errors.MediaSipPlayError.AddDetail("播放参数错误!")
	}
	return res, nil
}

func sipPlay(data *Stream) (*Stream, error) {
	chnRepo := db.NewSipChannelsRepo(Ctx)
	filter := db.SipChannelsFilter{ChannelIDs: []string{data.ChannelID}}
	chn, err := chnRepo.FindOneByFilter(Ctx, filter)
	if err != nil {
		fmt.Println("-----------------VidmgrGbsipChannelPlay  channelRepo Error----------------------")
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.MediaSipPlayError.AddDetail("通道不存在，ChnID:" + string(data.ChannelID))
		}
		return nil, err
	}
	chn.DeviceID = data.DeviceID

	// 推流模式要求设备在线且活跃
	if time.Now().Unix()-chn.LastLogin > 30*60 || chn.Status != DeviceStatusON {
		return nil, errors.MediaSipPlayError.AddDetail("通道已离线")
	}
	user, ok := activeDevices.Get(chn.DeviceID)
	if !ok {
		return nil, errors.MediaSipPlayError.AddDetail("设备已离线")
	}
	// GB28181推流
	if chn.Stream == "" {
		ssrcLock.Lock()
		data.ssrc = getSSRC(data.Type)
		chn.Stream = ssrc2stream(data.ssrc)
		// 成功后保存   更新chn的ID
		chnRepo.Update(Ctx, chn)
		ssrcLock.Unlock()
	}

	data, err = sipPlayPush(data, chn, &user)
	if err != nil {
		return nil, fmt.Errorf("获取视频失败:%v", err)
	}
	//data.HTTP = fmt.Sprintf("%s/rtp/%s/hls.m3u8", config.Media.HTTP, data.Stream)
	//data.RTMP = fmt.Sprintf("%s/rtp/%s", config.Media.RTMP, data.Stream)
	//data.RTSP = fmt.Sprintf("%s/rtp/%s", config.Media.RTSP, data.Stream)
	//data.WSFLV = fmt.Sprintf("%s/rtp/%s.live.flv", config.Media.WS, data.Stream)

	return data, nil
}

func sipPlayPush(data *Stream, chn *db.SipChannels, dev *db.SipDevices) (*Stream, error) {
	var (
		s sdp2.Session
		b []byte
	)
	name := "Play"
	protocal := "TCP/RTP/AVP"
	if data.Type == 1 { // 0  直播 1 历史
		name = "Playback"
		protocal = "RTP/RTCP"
	}

	video := sdp2.Media{
		Description: sdp2.MediaDescription{
			Type:     "video",
			Port:     int(dev.MediaPort),
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
	msg := &sdp2.Message{
		Origin: sdp2.Origin{
			Username: dev.DeviceID, // 媒体服务器id
			Address:  utils.InetNtoA(dev.MediaIP),
		},
		Name: name,
		Connection: sdp2.ConnectionData{
			IP:  net.ParseIP(utils.InetNtoA(dev.MediaIP)),
			TTL: 0,
		},
		Timing: []sdp2.Timing{
			{
				Start: data.S,
				End:   data.E,
			},
		},
		Medias: []sdp2.Media{video},
		SSRC:   data.ssrc,
	}
	if data.Type == 1 {
		msg.URI = fmt.Sprintf("%s:0", chn.ChannelID)
	}

	// appending message to session
	s = msg.Append(s)
	// appending session to byte buffer
	b = s.AppendTo(b)
	uri, _ := sip2.ParseURI(chn.URIStr)
	chn.Taddr = &sip2.Address{URI: uri}
	serverDevice.Taddr.Params.Add("tag", sip2.String{Str: utils.RandString(20)})
	hb := sip2.NewHeaderBuilder().SetTo(chn.Taddr).SetFrom(serverDevice.Taddr).AddVia(&sip2.ViaHop{
		Params: sip2.NewParams().Add("branch", sip2.String{Str: sip2.GenerateBranch()}),
	}).SetContentType(&sip2.ContentTypeSDP).SetMethod(sip2.INVITE).SetContact(serverDevice.Taddr)
	req := sip2.NewRequest("", sip2.INVITE, chn.Taddr.URI, sip2.DefaultSipVersion, hb.Build(), b)
	req.SetDestination(dev.Tsource)
	req.AppendHeader(&sip2.GenericHeader{HeaderName: "Subject", Contents: fmt.Sprintf("%s:%s,%s:%s", chn.ChannelID,
		chn.Stream, serverDevice.DeviceID, chn.Stream)})
	req.SetRecipient(chn.Taddr.URI)
	tx, err := SipSrv.Srv.Request(req)
	if err != nil {
		logrus.Warningln("sipPlayPush fail.id:", dev.DeviceID, chn.ChannelID, "err:", err)
		return data, err
	}
	// response
	response, err := sipResponse(tx)
	if err != nil {
		logrus.Warningln("sipPlayPush response fail.id:", dev.DeviceID, chn.ChannelID, "err:", err)
		return data, err
	}
	data.Resp = response
	// ACK
	tx.Request(sip2.NewRequestFromResponse(sip2.ACK, response))

	callid, _ := response.CallID()
	data.CallID = string(*callid)

	cseq, _ := response.CSeq()
	if cseq != nil {
		data.CseqNo = cseq.SeqNo
	}

	from, _ := response.From()
	to, _ := response.To()
	for k, v := range to.Params.Items() {
		data.Ttag[k] = v.String()
	}
	for k, v := range from.Params.Items() {
		data.Ftag[k] = v.String()
	}
	data.Status = 0

	return data, err
}

// 当前系统中存在的流列表
type streamsList struct {
	// key=ssrc value=PlayParams  播放对应的PlayParams 用来发送bye获取tag，callid等数据
	Response *sync.Map
	// key=channelid value={Play}  当前设备直播信息，防止重复直播
	Succ *sync.Map
	ssrc int
}

var StreamList streamsList

// 这个地方需要再研究一下
func getSSRC(t int) string {
	r := false
	for {
		StreamList.ssrc++
		// ssrc最大为四位数，超过时从1开始重新计算
		if StreamList.ssrc > 9000 && !r {
			StreamList.ssrc = 0
			r = true
		}
		key := fmt.Sprintf("%d%s%04d", t, SipInfo.Region[3:8], StreamList.ssrc)

		chnRepo := db.NewSipChannelsRepo(Ctx)
		filter := db.SipChannelsFilter{Stream: ssrc2stream(key)}
		_, err := chnRepo.FindOneByFilter(Ctx, filter)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return key
		}
	}
}

// zlm接收到的ssrc为16进制。发起请求的ssrc为10进制
func ssrc2stream(ssrc string) string {
	if ssrc[0:1] == "0" {
		ssrc = ssrc[1:]
	}
	num, _ := strconv.Atoi(ssrc)
	return fmt.Sprintf("%08X", num)
}
