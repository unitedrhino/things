package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/deviceinteract"
	server "github.com/i-Things/things/src/dmsvr/internal/server/deviceinteract"
)

var (
	deviceInteractSvr client.DeviceInteract
)

func NewDeviceInteract(runSvr bool) client.DeviceInteract {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectDeviceInteract(svcCtx, server.NewDeviceInteractServer(svcCtx))
	return dmSvr
}
