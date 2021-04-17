package svc

import (
	"github.com/tal-tech/go-zero/rest"
	"yl/src/webapi/api/test/internal/config"
	"yl/src/webapi/api/test/internal/middleware"
)

type ServiceContext struct {
	Config    config.Config
	Usercheck rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Usercheck: middleware.NewUsercheckMiddleware().Handle,
	}
}
