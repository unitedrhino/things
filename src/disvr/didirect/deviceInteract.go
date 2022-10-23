package didirect

import (
	client "github.com/i-Things/things/src/disvr/client/deviceinteract"
	server "github.com/i-Things/things/src/disvr/internal/server/deviceinteract"
)

var (
	deviceInteractSvr client.DeviceInteract
)

func NewDeviceInteract() client.DeviceInteract {
	svc := GetCtxSvc()
	dmSvr := client.NewDirectDeviceInteract(svc, server.NewDeviceInteractServer(svc))
	return dmSvr
}
