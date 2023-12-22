package dgdirect

import (
	client "github.com/i-Things/things/src/dgsvr/client/deviceauth"
	server "github.com/i-Things/things/src/dgsvr/internal/server/deviceauth"
)

var (
	deviceAuthSvr client.DeviceAuth
)

func NewDeviceAuth(runSvr bool) client.DeviceAuth {
	svcCtx := GetSvcCtx(runSvr)
	dgSvr := client.NewDirectDeviceAuth(svcCtx, server.NewDeviceAuthServer(svcCtx))
	return dgSvr
}
