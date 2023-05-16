package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/remoteconfig"
	server "github.com/i-Things/things/src/dmsvr/internal/server/remoteconfig"
)

var (
	remoteConfigSvr client.RemoteConfig
)

func NewRemoteConfig(runSvr bool) client.RemoteConfig {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	dmSvr := client.NewDirectRemoteConfig(svcCtx, server.NewRemoteConfigServer(svcCtx))
	return dmSvr
}
