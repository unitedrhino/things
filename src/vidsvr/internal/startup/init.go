package startup

import (
	//"github.com/i-Things/things/src/vidsvr/internal/event/busEvent/deviceDelete"
	//"github.com/i-Things/things/src/vidsvr/internal/event/busEvent/productSchemaUpdate"
	//"github.com/i-Things/things/src/vidsvr/internal/event/dataUpdateEvent"
	//"github.com/i-Things/things/src/vidsvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
)

func Init(svcCtx *svc.ServiceContext) {
	//InitSubscribe(svcCtx)
	InitEventBus(svcCtx)
}

//func InitSubscribe(svcCtx *svc.ServiceContext) {
//	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
//	if err != nil {
//		logx.Error("NewDataUpdate err", err)
//		os.Exit(-1)
//	}
//	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
//		return dataUpdateEvent.NewDataUpdateLogic(ctx, svcCtx)
//	})
//}

func InitEventBus(svcCtx *svc.ServiceContext) {
	//svcCtx.Bus.Subscribe(topics.DmDeviceInfoDelete, deviceDelete.DeviceGroupHandle(svcCtx))
	//svcCtx.Bus.Subscribe(topics.DmProductSchemaUpdate, productSchemaUpdate.EventsHandle(svcCtx))
}
