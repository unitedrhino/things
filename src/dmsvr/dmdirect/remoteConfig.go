package dmdirect

import (
	client "github.com/i-Things/things/src/dmsvr/client/remoteconfig"
	server "github.com/i-Things/things/src/dmsvr/internal/server/remoteconfig"
)

var (
	remoteConfigSvr client.RemoteConfig
)

func NewRemoteConfig() client.RemoteConfig {
	svc := GetCtxSvc()
	dmSvr := client.NewDirectRemoteConfig(svc, server.NewRemoteConfigServer(svc))
	return dmSvr
}
