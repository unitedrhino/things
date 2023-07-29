package startup

import (
	"context"
	"github.com/i-Things/things/src/stocksvr/internal/event/appDeviceEvent"
	"github.com/i-Things/things/src/stocksvr/internal/subscribe/subApp"
	"github.com/i-Things/things/src/stocksvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func Init(svcCtx *svc.ServiceContext) {
	Subscribe(svcCtx)
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
}
