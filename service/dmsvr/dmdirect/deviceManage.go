package dmdirect

import (
	client "gitee.com/i-Things/things/service/dmsvr/client/devicemanage"
	server "gitee.com/i-Things/things/service/dmsvr/internal/server/devicemanage"
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
