package dmdirect

import (
	client "github.com/i-Things/things/src/disvr/client/devicemsg"
	server "github.com/i-Things/things/src/disvr/internal/server/devicemsg"
)

var (
	deviceLogSvr client.DeviceMsg
)

func NewDeviceMsg(config *Config) client.DeviceMsg {
	svc := getCtxSvc(config)
	dmSvr := client.NewDirectDeviceMsg(svc, server.NewDeviceMsgServer(svc))
	return dmSvr
}
