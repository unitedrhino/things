package startup

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/eventBus"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/event/deviceMsgEvent"
	"github.com/i-Things/things/src/dmsvr/internal/event/serverEvent"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/subscribe/server"
	"github.com/i-Things/things/src/dmsvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {
	InitSubscribe(svcCtx)
	InitEventBus(svcCtx)
}
func InitSubscribe(svcCtx *svc.ServiceContext) {
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
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	svcCtx.ServerMsg.Subscribe(eventBus.DmDeviceInfoDelete, func(ctx context.Context, body []byte) error {
		var value devices.Core
		err := json.Unmarshal(body, &value)
		if err != nil {
			return err
		}
		err = relationDB.NewGroupDeviceRepo(ctx).DeleteByFilter(ctx, relationDB.GroupDeviceFilter{
			ProductID:  value.ProductID,
			DeviceName: value.DeviceName,
		})
		logx.WithContext(ctx).Infof("DeviceGroupHandle value:%v err:%v", utils.Fmt(value), err)
		return err
	})
	svcCtx.ServerMsg.Subscribe(eventBus.DmProductSchemaUpdate, func(ctx context.Context, body []byte) error {
		var productID = string(body)
		return svcCtx.SchemaRepo.ClearCache(ctx, productID)
	})
	err := svcCtx.ServerMsg.Start()
	logx.Must(err)
}
