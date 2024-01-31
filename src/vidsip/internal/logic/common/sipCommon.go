package common

import (
	"gitee.com/i-Things/core/shared/def"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/src/vidsip/pb/sip"
)

func ToPageInfo(info *sip.PageInfo, defaultOrders ...def.OrderBy) *def.PageInfo {
	if info == nil {
		return nil
	}

	var orders = defaultOrders
	if infoOrders := info.GetOrders(); len(infoOrders) > 0 {
		orders = make([]def.OrderBy, 0, len(infoOrders))
		for _, infoOd := range infoOrders {
			if infoOd.GetFiled() != "" {
				orders = append(orders, def.OrderBy{infoOd.GetFiled(), infoOd.GetSort()})
			}
		}
	}

	return &def.PageInfo{
		Page:   info.GetPage(),
		Size:   info.GetSize(),
		Orders: orders,
	}
}
func ToPageInfoWithDefault(info *sip.PageInfo, defau *def.PageInfo) *def.PageInfo {
	if page := ToPageInfo(info); page == nil {
		return defau
	} else {
		if page.Page == 0 {
			page.Page = defau.Page
		}
		if page.Size == 0 {
			page.Size = defau.Size
		}
		if len(page.Orders) == 0 {
			page.Orders = defau.Orders
		}
		return page
	}
}

/*
根据用户的输入生成对应的数据库数据
*/
func ToSipChannelDB(v *sip.SipChannel) (*db.SipChannels, error) {
	pi := &db.SipChannels{
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
	return pi, nil
}

func ToSipDeviceDB(v *sip.SipDevice) (*db.SipDevices, error) {
	// 调用 ResolveTCPAddr() 函数进行转换
	//natAddr, _ := net.ResolveTCPAddr("tcp", v.Source)
	pi := &db.SipDevices{
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
		LastLogin: v.LastLogin,
	}
	return pi, nil
}

/*
根据用户的输入生成对应的数据库数据
*/
func ToSipChannelRpc(v *db.SipChannels) *sip.SipChannel {
	pi := &sip.SipChannel{
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
		LastLogin:    v.LastLogin,
	}
	return pi
}

func ToSipDeviceRpc(v *db.SipDevices) *sip.SipDevice {
	pi := &sip.SipDevice{
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
		LastLogin: v.LastLogin,
	}
	return pi
}

func UpdatSipChannelDB(old *db.SipChannels, data *sip.SipChnUpdateReq) error {
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

func UpdatSipDeviceDB(old *db.SipDevices, data *sip.SipDevUpdateReq) error {
	if data.Name != "" {
		old.Name = data.Name
	}
	if data.PWD != "" {
		old.PWD = data.PWD
	}

	if data.VidmgrID != "" {
		old.VidmgrID = data.VidmgrID
	}
	return nil
}
