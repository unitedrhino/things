package svc

import (
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/publish/pubInner"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type ServiceContext struct {
	Config   config.Config
	PubDev   pubDev.PubDev
	PubInner pubInner.PubInner
}

func NewServiceContext(c config.Config) *ServiceContext {
	dl, err := pubDev.NewPubDev(c.DevLink)
	if err != nil {
		logx.Error("NewDevClient err", err)
		os.Exit(-1)
	}

	il, err := pubInner.NewPubInner(c.Event)
	if err != nil {
		logx.Error("NewInnerDevPub err", err)
		os.Exit(-1)
	}

	return &ServiceContext{
		Config:   c,
		PubDev:   dl,
		PubInner: il,
	}
}
