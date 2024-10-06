package dgdirect

import (
	client "gitee.com/i-Things/things/service/dgsvr/client/deviceauth"
	server "gitee.com/i-Things/things/service/dgsvr/internal/server/deviceauth"
)

var (
	deviceAuthSvr client.DeviceAuth
)

func NewDeviceAuth(runSvr bool) client.DeviceAuth {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dgSvr := client.NewDirectDeviceAuth(svcCtx, server.NewDeviceAuthServer(svcCtx))
	return dgSvr
}
