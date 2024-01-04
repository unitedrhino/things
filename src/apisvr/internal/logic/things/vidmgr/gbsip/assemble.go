package gbsip

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
)

func VidmgrGbsipDeviceToApi(v *vid.VidmgrGbsipDevice) *types.CommonSipDevice {
	return &types.CommonSipDevice{
		ID:           v.ID,
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
		LastLogin:    v.LastLogin,
		Regist:       v.Regist,
		PWD:          v.PWD,
		Source:       v.Source,
	}
}

func VidmgrGbsipChanneloApi(v *vid.VidmgrGbsipChannel) *types.CommonSipChannel {
	return &types.CommonSipChannel{
		ID:           v.ID,
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
		LastLogin:    v.LastLogin,
	}
}
