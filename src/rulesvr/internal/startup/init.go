package startup

import (
	"context"
	"gitee.com/i-Things/core/shared/events/topics"
	"github.com/i-Things/things/src/rulesvr/internal/event/appDeviceEvent"
	"github.com/i-Things/things/src/rulesvr/internal/event/busEvent/sceneChange"
	"github.com/i-Things/things/src/rulesvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/rulesvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/rulesvr/internal/repo/event/subscribe/subApp"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {
	Subscribe(svcCtx)
	InitEventBus(svcCtx)
}

func Subscribe(svcCtx *svc.ServiceContext) {
	subAppCli, err := subApp.NewSubApp(svcCtx.Config.Event)
	logx.Must(err)
	err = subAppCli.Subscribe(func(ctx context.Context) subApp.AppSubEvent {
		return appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx)
	})
	logx.Must(err)
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	logx.Must(err)
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	})
	logx.Must(err)
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	svcCtx.Bus.Subscribe(topics.RuleSceneInfoDelete, sceneChange.EventsHandle(svcCtx, topics.RuleSceneInfoDelete))
	svcCtx.Bus.Subscribe(topics.RuleSceneInfoUpdate, sceneChange.EventsHandle(svcCtx, topics.RuleSceneInfoUpdate))
}
