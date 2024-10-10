package dmdirect

import (
	client "gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	server "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemanage"
)

var (
	deviceManageSvr client.DeviceManage
)

func NewDeviceManage(runSvr bool) client.DeviceManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectDeviceManage(svcCtx, server.NewDeviceManageServer(svcCtx))
	return dmSvr
}
