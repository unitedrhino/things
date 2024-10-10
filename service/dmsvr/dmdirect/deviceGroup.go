package dmdirect

import (
	client "gitee.com/unitedrhino/things/service/dmsvr/client/devicegroup"
	server "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicegroup"
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
