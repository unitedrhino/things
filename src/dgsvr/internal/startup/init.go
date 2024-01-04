package startup

import (
	"context"
	"github.com/i-Things/things/src/dgsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/dgsvr/internal/event/deviceSub"
	"github.com/i-Things/things/src/dgsvr/internal/event/innerSub"
	"github.com/i-Things/things/src/dgsvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/dgsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/src/dgsvr/internal/repo/event/publish/pubInner"
	"github.com/i-Things/things/src/dgsvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/src/dgsvr/internal/repo/event/subscribe/subInner"
	"github.com/i-Things/things/src/dgsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {
	//some init for serviceContext
}

// mqtt and nats client
func PostInit(svcCtx *svc.ServiceContext) {
	dl, err := pubDev.NewPubDev(svcCtx.Config.DevLink)
	logx.Must(err)

	il, err := pubInner.NewPubInner(svcCtx.Config.Event)
	logx.Must(err)

	svcCtx.PubDev = dl
	svcCtx.PubInner = il

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
