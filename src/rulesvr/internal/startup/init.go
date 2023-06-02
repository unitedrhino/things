package startup

import (
	"context"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/event/appDeviceEvent"
	"github.com/i-Things/things/src/rulesvr/internal/event/busEvent/sceneChange"
	"github.com/i-Things/things/src/rulesvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/rulesvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/rulesvr/internal/repo/event/subscribe/subApp"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"os"
)

func Init(svcCtx *svc.ServiceContext) {
	Subscribe(svcCtx)
	InitEventBus(svcCtx)
}

func Subscribe(svcCtx *svc.ServiceContext) {
	subAppCli, err := subApp.NewSubApp(svcCtx.Config.Event)
	if err != nil {
		logx.Error("NewSubApp err", err)
		os.Exit(-1)
	}
	err = subAppCli.Subscribe(func(ctx context.Context) subApp.AppSubEvent {
		return appDeviceEvent.NewAppDeviceHandle(ctx, svcCtx)
	})
	if err != nil {
		log.Fatalf("%v.subApp.Subscribe err:%v",
			utils.FuncName(), err)
	}
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	})
	if err != nil {
		log.Fatalf("%v.DataUpdate err:%v",
			utils.FuncName(), err)
	}
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	svcCtx.Bus.Subscribe(topics.RuleSceneInfoDelete, sceneChange.EventsHandle(svcCtx, topics.RuleSceneInfoDelete))
	svcCtx.Bus.Subscribe(topics.RuleSceneInfoUpdate, sceneChange.EventsHandle(svcCtx, topics.RuleSceneInfoUpdate))
}
