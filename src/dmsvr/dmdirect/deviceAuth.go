package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/deviceauth"
	server "github.com/i-Things/things/src/dmsvr/internal/server/deviceauth"
)

var (
	deviceAuthSvr client.DeviceAuth
)

func NewDeviceAuth() client.DeviceAuth {
	svc := GetCtxSvc()
	dmSvr := client.NewDirectDeviceAuth(svc, server.NewDeviceAuthServer(svc))
	return dmSvr
}
