package startup

import (
	"context"
	"github.com/i-Things/things/src/disvr/internal/event/dataUpdateEvent"
	"github.com/i-Things/things/src/disvr/internal/event/deviceMsgEvent"
	"github.com/i-Things/things/src/disvr/internal/event/serverEvent"
	"github.com/i-Things/things/src/disvr/internal/repo/event/subscribe/dataUpdate"
	"github.com/i-Things/things/src/disvr/internal/repo/event/subscribe/server"
	"github.com/i-Things/things/src/disvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Subscribe(svcCtx *svc.ServiceContext) {
	{
		cli, err := subDev.NewSubDev(svcCtx.Config.Event)
		logx.Must(err)
		err = cli.Subscribe(func(ctx context.Context) subDev.InnerSubEvent {
			return deviceMsgEvent.NewDeviceMsgHandle(ctx, svcCtx)
		})
		logx.Must(err)
	}
	{
		cli, err := server.NewServer(svcCtx.Config.Event)
		logx.Must(err)
		err = cli.Subscribe(func(ctx context.Context) server.ServerHandle {
			return serverEvent.NewServerHandle(ctx, svcCtx)
		})
		logx.Must(err)
	}
	{
		cli, err := dataUpdate.NewDataUpdate(svcCtx.Config.Event)
		logx.Must(err)
		err = cli.Subscribe(func(ctx context.Context) dataUpdate.UpdateHandle {
			return dataUpdateEvent.NewPublishLogic(ctx, svcCtx)
		})
		logx.Must(err)
	}
}
