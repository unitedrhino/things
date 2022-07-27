package svc

import (
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/devLink"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/innerLink"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type ServiceContext struct {
	DevLink   devLink.DevLink
	InnerLink innerLink.InnerLink
}

func NewServiceContext(c config.Config) *ServiceContext {
	dl, err := devLink.NewDevClient(c.DevLink)
	if err != nil {
		logx.Error("NewDevClient err", err)
		os.Exit(-1)
	}
	il, err := innerLink.NewInnerLink(c.InnerLink)
	if err != nil {
		logx.Error("NewInnerLink err", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		DevLink:   dl,
		InnerLink: il,
	}
}
