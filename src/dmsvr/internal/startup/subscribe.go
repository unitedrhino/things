package startup

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func Subscribe(svcCtx *svc.ServiceContext) {
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	})
}
