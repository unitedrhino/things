package didirect

import (
	client "github.com/i-Things/things/src/disvr/client/devicemsg"
	server "github.com/i-Things/things/src/disvr/internal/server/devicemsg"
)

var (
	deviceLogSvr client.DeviceMsg
)

func NewDeviceMsg() client.DeviceMsg {
	svc := GetCtxSvc()
	dmSvr := client.NewDirectDeviceMsg(svc, server.NewDeviceMsgServer(svc))
	return dmSvr
}
