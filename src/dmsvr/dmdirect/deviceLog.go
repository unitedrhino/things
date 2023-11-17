package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/devicemsg"
	server "github.com/i-Things/things/src/dmsvr/internal/server/devicemsg"
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
