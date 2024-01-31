package svc

import (
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/viewsvr/internal/config"
	"github.com/i-Things/things/service/viewsvr/internal/middleware"
	"github.com/i-Things/things/service/viewsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"os"
)

type ServiceContext struct {
	Config         config.Config
	SetupWare      rest.Middleware
	CheckTokenWare rest.Middleware
	DataAuthWare   rest.Middleware
	TeardownWare   rest.Middleware
	CheckApiWare   rest.Middleware
	OssClient      *oss.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	relationDB.Migrate(c.Database)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:         c,
		SetupWare:      middleware.NewSetupWareMiddleware().Handle,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware().Handle,
		DataAuthWare:   middleware.NewDataAuthWareMiddleware().Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware().Handle,
		CheckApiWare:   middleware.NewCheckApiWareMiddleware().Handle,
		OssClient:      ossClient,
	}
}
