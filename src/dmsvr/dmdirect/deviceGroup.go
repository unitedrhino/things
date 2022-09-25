package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/devicegroup"
	server "github.com/i-Things/things/src/dmsvr/internal/server/devicegroup"
)

var (
	deviceGroupSvr client.DeviceGroup
)

func NewDeviceGroup(config *Config) client.DeviceGroup {
	svc := getCtxSvc(config)
	dmSvr := client.NewDirectDeviceGroup(svc, server.NewDeviceGroupServer(svc))
	return dmSvr
}
