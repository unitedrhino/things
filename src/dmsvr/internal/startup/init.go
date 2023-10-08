package startup

import (
	"context"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/src/dmsvr/internal/event/busEvent/deviceDelete"
	"github.com/i-Things/things/src/dmsvr/internal/event/busEvent/productSchemaUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {
	InitSubscribe(svcCtx)
	InitEventBus(svcCtx)
}
func InitSubscribe(svcCtx *svc.ServiceContext) {
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	logx.Must(err)
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewDataUpdateLogic(ctx, svcCtx)
	})
	logx.Must(err)
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	svcCtx.Bus.Subscribe(topics.DmDeviceInfoDelete, deviceDelete.DeviceGroupHandle(svcCtx))
	svcCtx.Bus.Subscribe(topics.DmProductSchemaUpdate, productSchemaUpdate.EventsHandle(svcCtx))
}
