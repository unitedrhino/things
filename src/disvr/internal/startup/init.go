package startup

import (
	"context"
	"github.com/i-Things/things/src/disvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/disvr/internal/event/deviceMsgEvent"
	"github.com/i-Things/things/src/disvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/disvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Subscribe(svcCtx *svc.ServiceContext) {
	subDevCli, err := subDev.NewSubDev(svcCtx.Config.Event)
	logx.Must(err)
	err = subDevCli.Subscribe(func(ctx context.Context) subDev.InnerSubEvent {
		return deviceMsgEvent.NewDeviceMsgHandle(ctx, svcCtx)
	})
	logx.Must(err)
	dataUpdateCli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
	logx.Must(err)
	err = dataUpdateCli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
		return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
	})
	logx.Must(err)
}
