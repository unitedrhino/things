package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/devicegroup"
	server "github.com/i-Things/things/src/dmsvr/internal/server/devicegroup"
)

var (
	deviceGroupSvr client.DeviceGroup
)

func NewDeviceGroup(runSvr bool) client.DeviceGroup {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectDeviceGroup(svcCtx, server.NewDeviceGroupServer(svcCtx))
	return dmSvr
}
