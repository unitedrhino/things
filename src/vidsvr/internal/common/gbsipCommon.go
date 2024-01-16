package common

import (
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"time"
)

/*
根据用户的输入生成对应的数据库数据
*/
func ToVidmgrGbsipChannelDB(v *vid.VidmgrGbsipChannel) (*relationDB.VidmgrChannels, error) {
	pi := &relationDB.VidmgrChannels{
		ChannelID:    v.ChannelID,
		DeviceID:     v.DeviceID,
		Memo:         v.Memo,
		Name:         v.Name,
		Manufacturer: v.Manufacturer,
		Model:        v.Model,
		Owner:        v.Owner,
		CivilCode:    v.CivilCode,
		Address:      v.Address,
		Parental:     v.Parental,
		SafetyWay:    v.SafetyWay,
		RegisterWay:  v.RegisterWay,
		Secrecy:      v.Secrecy,
		Status:       v.Status,
		URIStr:       v.URIStr,
		VF:           v.VF,
		Height:       v.Height,
		Width:        v.Width,
		FPS:          v.FPS,
		StreamType:   v.StreamType,
		URL:          v.URL,
		LastLogin:    time.Unix(v.LastLogin, 0),
		IsPlay:       v.IsPlay,
	}
	return pi, nil
}

func ToVidmgrGbsipDeviceDB(v *vid.VidmgrGbsipDevice) (*relationDB.VidmgrDevices, error) {
	// 调用 ResolveTCPAddr() 函数进行转换
	//natAddr, _ := net.ResolveTCPAddr("tcp", v.Source)
	pi := &relationDB.VidmgrDevices{
		DeviceID:     v.DeviceID,
		Name:         v.Name,
		Region:       v.Region,
		Host:         v.Host,
		Port:         v.Port,
		TransPort:    v.TransPort,
		Proto:        v.Proto,
		Rport:        v.Rport,
		RAddr:        v.RAddr,
		Manufacturer: v.Manufacturer,
		DeviceType:   v.DeviceType,
		Firmware:     v.Firmware,
		Model:        v.Model,
		URIStr:       v.URIStr,
		Regist:       v.Regist,
		PWD:          v.PWD,
		//Source:       natAddr,
		LastLogin: time.Unix(v.LastLogin, 0),
	}
	return pi, nil
}

/*
根据用户的输入生成对应的数据库数据
*/
func ToVidmgrGbsipChannelRpc(v *relationDB.VidmgrChannels) *vid.VidmgrGbsipChannel {
	pi := &vid.VidmgrGbsipChannel{
		ChannelID:    v.ChannelID,
		DeviceID:     v.DeviceID,
		Memo:         v.Memo,
		Name:         v.Name,
		Manufacturer: v.Manufacturer,
		Model:        v.Model,
		Owner:        v.Owner,
		CivilCode:    v.CivilCode,
		Address:      v.Address,
		Parental:     v.Parental,
		SafetyWay:    v.SafetyWay,
		RegisterWay:  v.RegisterWay,
		Secrecy:      v.Secrecy,
		Status:       v.Status,
		URIStr:       v.URIStr,
		VF:           v.VF,
		IsPlay:       v.IsPlay,
		Height:       v.Height,
		Width:        v.Width,
		FPS:          v.FPS,
		StreamType:   v.StreamType,
		URL:          v.URL,
		LastLogin:    v.LastLogin.Unix(),
	}
	return pi
}

func ToVidmgrGbsipDeviceRpc(v *relationDB.VidmgrDevices) *vid.VidmgrGbsipDevice {
	pi := &vid.VidmgrGbsipDevice{
		DeviceID:     v.DeviceID,
		Name:         v.Name,
		Region:       v.Region,
		Host:         v.Host,
		Port:         v.Port,
		TransPort:    v.TransPort,
		Proto:        v.Proto,
		Rport:        v.Rport,
		RAddr:        v.RAddr,
		Manufacturer: v.Manufacturer,
		DeviceType:   v.DeviceType,
		Firmware:     v.Firmware,
		Model:        v.Model,
		URIStr:       v.URIStr,
		Regist:       v.Regist,
		PWD:          v.PWD,
		//Source:       v.Source.String(),
		LastLogin: v.LastLogin.Unix(),
	}
	return pi
}

func UpdatVidmgrChannelDB(old *relationDB.VidmgrChannels, data *vid.VidmgrGbsipChannelUpdate) error {
	if data.Memo != "" {
		old.Name = data.Memo
	}
	if data.StreamType != "" {
		old.StreamType = data.StreamType
	}
	if data.Url != "" {
		old.URL = data.Url
	}

	return nil
}

func UpdatVidmgrDeviceDB(old *relationDB.VidmgrDevices, data *vid.VidmgrGbsipDeviceUpdateReq) error {
	if data.Name != "" {
		old.Name = data.Name
	}
	if data.PWD != "" {
		old.PWD = data.PWD
	}
	return nil
}
