package dmdirect

import (
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/server"
	"github.com/i-Things/things/src/dmsvr/internal/startup"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
)

type Config = config.Config

var (
	ctxSvc *svc.ServiceContext
	dmSvr  dm.Dm
)

func NewDm(config *Config) dm.Dm {
	if dmSvr != nil {
		return dmSvr
	}
	dmSvc := svc.NewServiceContext(*config)
	startup.Subscribe(dmSvc)
	dmSvr = dm.NewDirectDm(dmSvc, server.NewDmServer(dmSvc))
	return dmSvr
}
