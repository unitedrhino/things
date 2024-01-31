package svc

import (
	"gitee.com/i-Things/core/shared/stores"
	"github.com/i-Things/things/service/udsvr/internal/config"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
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
