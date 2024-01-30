package gbsip

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/vidsip/pb/sip"
)

func VidmgrGbsipDeviceToApi(v *sip.SipDevice) *types.CommonSipDevice {
	return &types.CommonSipDevice{
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

func VidmgrGbsipChanneloApi(v *sip.SipChannel) *types.CommonSipChannel {
	return &types.CommonSipChannel{
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
		IsPlay:       v.IsPlay,
	}
}

func ToVidmgrGbsipInfoApi(v *sip.SipInfo) *types.CommonSipInfo {
	return &types.CommonSipInfo{
		VidmgrID:     v.VidmgrID,
		ID:           v.ID,
		Region:       v.Region,
		CID:          v.CID,
		CNUM:         v.CNUM,
		DID:          v.DID,
		DNUM:         v.DNUM,
		LID:          v.LID,
		IsOpen:       v.IsOpen,
		IP:           v.IP,
		Port:         v.Port,
		MediaRtpIP:   utils.InetNtoA(v.MediaRtpIP),
		MediaRtpPort: v.MediaRtpPort,
	}
}
