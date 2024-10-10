package dmdirect

import (
	client "gitee.com/unitedrhino/things/service/dmsvr/client/userdevice"
	server "gitee.com/unitedrhino/things/service/dmsvr/internal/server/userdevice"
)

var (
	userDeviceSvr client.UserDevice
)

func NewUserDevice(runSvr bool) client.UserDevice {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectUserDevice(svcCtx, server.NewUserDeviceServer(svcCtx))
	return dmSvr
}
