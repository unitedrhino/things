package dgdirect

import (
	client "github.com/i-Things/things/src/dgsvr/client/deviceauth"
	server "github.com/i-Things/things/src/dgsvr/internal/server/deviceauth"
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
