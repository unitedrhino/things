package didirect

import (
	client "github.com/i-Things/things/src/disvr/client/deviceinteract"
	server "github.com/i-Things/things/src/disvr/internal/server/deviceinteract"
)

var (
	deviceInteractSvr client.DeviceInteract
)

func NewDeviceInteract(config *Config) client.DeviceInteract {
	svc := getCtxSvc(config)
	dmSvr := client.NewDirectDeviceInteract(svc, server.NewDeviceInteractServer(svc))
	return dmSvr
}
