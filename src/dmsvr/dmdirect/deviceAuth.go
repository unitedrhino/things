package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/deviceauth"
	server "github.com/i-Things/things/src/dmsvr/internal/server/deviceauth"
)

var (
	deviceAuthSvr client.DeviceAuth
)

func NewDeviceAuth(runSvr bool) client.DeviceAuth {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectDeviceAuth(svcCtx, server.NewDeviceAuthServer(svcCtx))
	return dmSvr
}
