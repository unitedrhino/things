package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/devicemanage"
	server "github.com/i-Things/things/src/dmsvr/internal/server/devicemanage"
)

var (
	deviceManageSvr client.DeviceManage
)

func NewDeviceManage() client.DeviceManage {
	svc := GetCtxSvc()
	dmSvr := client.NewDirectDeviceManage(svc, server.NewDeviceManageServer(svc))
	return dmSvr
}
