package startup

import (
	"context"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/src/dmsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/dmsvr/internal/event/eventChange/deviceDelete"
	"github.com/i-Things/things/src/dmsvr/internal/event/eventChange/productSchemaUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

func Init(svcCtx *svc.ServiceContext) {
	InitSubscribe(svcCtx)
	InitEventBus(svcCtx)
}
func InitSubscribe(svcCtx *svc.ServiceContext) {
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	if err != nil {
		logx.Error("NewDataUpdate err", err)
		os.Exit(-1)
	}
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewDataUpdateLogic(ctx, svcCtx)
	})
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	svcCtx.Bus.Subscribe(topics.DmDeviceDelete, deviceDelete.DeviceGroupHandle(svcCtx))
	svcCtx.Bus.Subscribe(topics.DmProductUpdateSchema, productSchemaUpdate.EventsHandle(svcCtx))
}
