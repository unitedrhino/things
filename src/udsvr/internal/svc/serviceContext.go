package svc

import (
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/udsvr/internal/config"
	"github.com/i-Things/things/src/udsvr/internal/repo/relationDB"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	relationDB.Migrate(c.Database)
	return &ServiceContext{
		Config: c,
	}
}
