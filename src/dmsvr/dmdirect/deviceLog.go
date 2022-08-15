package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/devicelog"
	server "github.com/i-Things/things/src/dmsvr/internal/server/devicelog"
)

var (
	deviceLogSvr client.DeviceLog
)

func NewDeviceLog(config *Config) client.DeviceLog {
	svc := getCtxSvc(config)
	dmSvr := client.NewDirectDeviceLog(svc, server.NewDeviceLogServer(svc))
	return dmSvr
}
