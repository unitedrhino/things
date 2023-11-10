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
	logx.Must(err)
	err = sd.SubDevMsg(func(ctx context.Context) subDev.DevSubHandle {
		return deviceSub.NewDeviceSubServer(svcCtx, ctx)
	})
	logx.Must(err)

	si, err := subInner.NewSubInner(svcCtx.Config.Event)
	logx.Must(err)
	err = si.SubToDevMsg(func(ctx context.Context) subInner.InnerSubHandle {
		return innerSub.NewInnerSubServer(svcCtx, ctx)
	})
	logx.Must(err)
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	logx.Must(err)
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewDataUpdateLogic(ctx, svcCtx)
	})
	logx.Must(err)
}
