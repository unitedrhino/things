package media

import (
	"bytes"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/gosip/sip"
	db "github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// 从请求中解析出设备信息
func parserDevicesFromReqeust(req *sip.Request) (db.VidmgrDevices, bool) {
	u := db.VidmgrDevices{}
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
	u.Addr = sip.NewAddressFromFromHeader(header)
	u.SourceStr = req.Source().String()
	u.Source = req.Source()
	return u, true
}

// 获取设备信息（注册设备）
func sipDeviceInfo(to db.VidmgrDevices) {
	hb := sip.NewHeaderBuilder().SetTo(to.Addr).SetFrom(_serverDevices.Addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE)
	req := sip.NewRequest("", sip.MESSAGE, to.Addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetDeviceInfoXML(to.DeviceID))
	req.SetDestination(to.Source)
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

// zlm接收到的ssrc为16进制。发起请求的ssrc为10进制
func ssrc2stream(ssrc string) string {
	if ssrc[0:1] == "0" {
		ssrc = ssrc[1:]
	}
	num, _ := strconv.Atoi(ssrc)
	return fmt.Sprintf("%08X", num)
}

func sipMessageCatalog(u db.VidmgrDevices, body []byte) error {
	message := &MessageDeviceListResponse{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	channelsRepo := db.NewVidmgrChannelsRepo(Ctx)
	if message.SumNum > 0 {
		for _, d := range message.Item {
			filter := db.VidmgrChannelsFilter{
				DeviceIDs:  []string{d.DeviceID},
				ChannelIDs: []string{d.ChannelID},
			}
			channel, err := channelsRepo.FindOneByFilter(Ctx, filter)
			if err == nil {
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
				go notify(notifyChannelsActive(*channel))
			} else {
				logrus.Infoln("deviceid not found,deviceid:", d.DeviceID, "pdid:", message.DeviceID, "err", err)
			}
		}
	}
	return nil
}
func sipMessageKeepalive(u db.VidmgrDevices, body []byte) error {
	message := &MessageNotify{}
	if err := utils.XMLDecode(body, message); err != nil {
		logrus.Errorln("Message Unmarshal xml err:", err, "body:", string(body))
		return err
	}
	device, ok := _activeDevices.Get(u.DeviceID)
	deviceRepo := db.NewVidmgrDevicesRepo(Ctx)
	if !ok {
		//device = db.VidmgrDevices{DeviceID: u.DeviceID}
		filter := db.VidmgrDevicesFilter{
			DeviceIDs: []string{u.DeviceID},
		}
		device1, err := deviceRepo.FindOneByFilter(Ctx, filter)
		if err != nil {
			logrus.Warnln("Device Keepalive not found ", u.DeviceID, err)
		}
		device = *device1
	}
	if message.Status == "OK" {
		device.LastLogin = time.Now()
		_activeDevices.Store(u.DeviceID, u)
	} else {
		//device.LastLogin = -1
		_activeDevices.Delete(u.DeviceID)
	}
	go notify(notifyDevicesAcitve(u.DeviceID, message.Status))
	device.Host = u.Host
	device.Port = u.Port
	device.Rport = u.Rport
	device.RAddr = u.RAddr
	device.Source = u.Source
	device.URIStr = u.URIStr

	err := deviceRepo.Update(Ctx, &device)

	return err
}

func sipMessageRecordInfo(u db.VidmgrDevices, body []byte) error {
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

func sipMessageDeviceInfo(u db.VidmgrDevices, body []byte) error {
	message := &MessageDeviceInfoResponse{}
	if err := utils.XMLDecode([]byte(body), message); err != nil {
		logrus.Errorln("sipMessageDeviceInfo Unmarshal xml err:", err, "body:", body)
		return err
	}
	deviceRepo := db.NewVidmgrDevicesRepo(Ctx)
	u.Model = message.Model
	u.DeviceType = message.DeviceType
	u.Firmware = message.Firmware
	u.Manufacturer = message.Manufacturer
	deviceRepo.Update(Ctx, &u)
	return nil
}

// sipCatalog 获取注册设备包含的列表
func sipCatalog(to db.VidmgrDevices) {
	hb := sip.NewHeaderBuilder().SetTo(to.Addr).SetFrom(_serverDevices.Addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE)
	req := sip.NewRequest("", sip.MESSAGE, to.Addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetCatalogXML(to.DeviceID))
	req.SetDestination(to.Source)
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
