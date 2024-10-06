package dmdirect

import (
	client "gitee.com/i-Things/things/service/dmsvr/client/remoteconfig"
	server "gitee.com/i-Things/things/service/dmsvr/internal/server/remoteconfig"
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
