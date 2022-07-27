package startup

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/dmsvr/internal/event/deviceMsgEvent"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/dataUpdate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/innerLink"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"log"
)

func Subscribe(svcCtx *svc.ServiceContext) {
	err := svcCtx.InnerLink.Subscribe(func(ctx context.Context) innerLink.InnerSubEvent {
		return deviceMsgEvent.NewDeviceMsgHandle(ctx, svcCtx)
	})
	if err != nil {
		log.Fatalf("%v|InnerLink.Subscribe|err:%v",
			utils.FuncName(), err)
	}
	err = svcCtx.DataUpdate.Subscribe(func(ctx context.Context) dataUpdate.DataUpdateSubEvent {
		return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	})
	if err != nil {
		log.Fatalf("[%v]DataUpdate.Subscribe|err:%v",
			utils.FuncName(), err)
	}
}
