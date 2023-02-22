package startup

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/event/appDeviceEvent"
	"github.com/i-Things/things/src/rulesvr/internal/repo/event/subscribe/subApp"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"os"
)

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
}
