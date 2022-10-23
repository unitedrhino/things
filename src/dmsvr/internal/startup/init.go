package startup

import (
	"github.com/i-Things/things/src/dmsvr/internal/svc"
)

func Subscribe(svcCtx *svc.ServiceContext) {
	//subDevCli, err := subDev.NewSubDev(svcCtx.Config.Event)
	//if err != nil {
	//	logx.Error("NewSubDev err", err)
	//	os.Exit(-1)
	//}
	//err = subDevCli.Subscribe(func(ctx context.Context) subDev.InnerSubEvent {
	//	return deviceMsgEvent.NewDeviceMsgHandle(ctx, svcCtx)
	//})
	//if err != nil {
	//	log.Fatalf("%v|SubDev.Subscribe|err:%v",
	//		utils.FuncName(), err)
	//}
	//dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	//if err != nil {
	//	logx.Error("NewDataUpdate err", err)
	//	os.Exit(-1)
	//}
	//err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.DataUpdateHandle {
	//	return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	//})
	//if err != nil {
	//	log.Fatalf("[%v]DataUpdate.Subscribe|err:%v",
	//		utils.FuncName(), err)
	//}
}
