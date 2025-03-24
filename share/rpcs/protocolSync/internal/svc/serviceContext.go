package svc

import "gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
