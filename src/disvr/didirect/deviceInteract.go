package didirect

import (
	client "github.com/i-Things/things/src/disvr/client/deviceinteract"
	server "github.com/i-Things/things/src/disvr/internal/server/deviceinteract"
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
