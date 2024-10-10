package dmdirect

import (
	client "gitee.com/unitedrhino/things/service/dmsvr/client/devicemsg"
	server "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemsg"
)

var (
	deviceLogSvr client.DeviceMsg
)

func NewDeviceMsg(runSvr bool) client.DeviceMsg {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectDeviceMsg(svcCtx, server.NewDeviceMsgServer(svcCtx))
	return dmSvr
}
