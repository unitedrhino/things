package startup

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/events"
	"github.com/i-Things/things/service/dgsvr/internal/event/deviceSub"
	"github.com/i-Things/things/service/dgsvr/internal/event/innerSub"
	"github.com/i-Things/things/service/dgsvr/internal/repo/event/publish/pubDev"
	"github.com/i-Things/things/service/dgsvr/internal/repo/event/publish/pubInner"
	"github.com/i-Things/things/service/dgsvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/service/dgsvr/internal/repo/event/subscribe/subInner"
	"github.com/i-Things/things/service/dgsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {
	//some init for serviceContext
}

// mqtt and nats client
func PostInit(svcCtx *svc.ServiceContext) {
	dl, err := pubDev.NewPubDev(svcCtx.Config.DevLink)
	logx.Must(err)

	il, err := pubInner.NewPubInner(svcCtx.Config.Event, def.ProtocolCodeIThings)
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
	InitEventBus(svcCtx)
}
func InitEventBus(svcCtx *svc.ServiceContext) {
	err := svcCtx.FastEvent.Subscribe(eventBus.DmProductCustomUpdate, func(ctx context.Context, t time.Time, body []byte) error {
		info := events.DeviceUpdateInfo{}
		err := json.Unmarshal(body, &info)
		if err != nil {
			return err
		}
		return svcCtx.Script.ClearCache(ctx, info.ProductID)
	})
	logx.Must(err)
	err = svcCtx.FastEvent.Start()
	logx.Must(err)
}
