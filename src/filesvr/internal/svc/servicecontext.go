package svc

import (
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/src/filesvr/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type ServiceContext struct {
	Config config.Config
	OSS    oss.OSSer
}

func NewServiceContext(c config.Config) *ServiceContext {
	ossClient, err := oss.NewOss(c.OSS)
	if err != nil {
		logx.Error("NewOss err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config: c,
		OSS:    ossClient,
	}
}
