package startup

import (
	"context"
	"fmt"
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/ddsvr/internal/event/deviceSub"
	"github.com/i-Things/things/src/ddsvr/internal/event/innerSub"
	"github.com/i-Things/things/src/ddsvr/internal/handler"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/src/ddsvr/internal/repo/event/subscribe/subInner"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"os"
)

func NewDd(c config.Config) {
	svcCtx := svc.NewServiceContext(c)
	Subscript(svcCtx)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	handler.RegisterHandlers(server, svcCtx)

	fmt.Printf("Starting ddSvr at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func Subscript(svcCtx *svc.ServiceContext) {
	sd, err := subDev.NewSubDev(svcCtx.Config.DevLink)
	if err != nil {
		logx.Error("NewSubDev err", err)
		os.Exit(-1)
	}
	err = sd.SubDevMsg(func(ctx context.Context) subDev.DevSubHandle {
		return deviceSub.NewDeviceSubServer(svcCtx, ctx)
	})
	if err != nil {
		logx.Error("PubDev SubToDevMsg failure", err)
		os.Exit(-1)
	}

	si, err := subInner.NewSubInner(svcCtx.Config.Event)
	if err != nil {
		logx.Error("NewSubInner err", err)
		os.Exit(-1)
	}
	err = si.SubToDevMsg(func(ctx context.Context) subInner.InnerSubHandle {
		return innerSub.NewInnerSubServer(svcCtx, ctx)
	})
	if err != nil {
		logx.Error("DevPubInner SubToDevMsg failure", err)
		os.Exit(-1)
	}
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewDataUpdateLogic(ctx, svcCtx)
	})
}
