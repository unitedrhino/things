package dcdirect

import (
	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/i-Things/things/src/dcsvr/internal/config"
	"github.com/i-Things/things/src/dcsvr/internal/server"
	"github.com/i-Things/things/src/dcsvr/internal/svc"
)

type Config = config.Config

var ctxSvc *svc.ServiceContext

func NewDc(config Config) dc.Dc {
	dcSvc := svc.NewServiceContext(config)
	return dc.NewDirectDc(dcSvc, server.NewDcServer(dcSvc))
}
