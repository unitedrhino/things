package startup

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/event/deviceMsgEvent"
	"github.com/i-Things/things/service/dmsvr/internal/event/serverEvent"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/event/subscribe/server"
	"github.com/i-Things/things/service/dmsvr/internal/repo/event/subscribe/subDev"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {
	InitCache(svcCtx)
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

func InitCache(svcCtx *svc.ServiceContext) {
	productCache, err := caches.NewCache(caches.CacheConfig[dm.ProductInfo]{
		KeyType:   eventBus.ServerCacheKeyDmProduct,
		FastEvent: svcCtx.ServerMsg,
		GetData: func(ctx context.Context, key string) (*dm.ProductInfo, error) {
			db := relationDB.NewProductInfoRepo(ctx)
			pi, err := db.FindOneByFilter(ctx, relationDB.ProductFilter{
				ProductIDs: []string{key}, WithProtocol: true, WithCategory: true})
			pb := logic.ToProductInfo(ctx, svcCtx, pi)
			return pb, err
		},
		ExpireTime: 10 * time.Minute,
	})
	logx.Must(err)
	svcCtx.ProductCache = productCache
	deviceCache, err := caches.NewCache(caches.CacheConfig[dm.DeviceInfo]{
		KeyType:   eventBus.ServerCacheKeyDmDevice,
		FastEvent: svcCtx.ServerMsg,
		GetData: func(ctx context.Context, key string) (*dm.DeviceInfo, error) {
			db := relationDB.NewDeviceInfoRepo(ctx)
			productID, deviceName, _ := strings.Cut(key, ":")
			di, err := db.FindOneByFilter(ctx, relationDB.DeviceFilter{
				ProductID: productID, DeviceName: deviceName})
			pb := logic.ToDeviceInfo(di)
			return pb, err
		},
		ExpireTime: 10 * time.Minute,
	})
	logx.Must(err)
	svcCtx.DeviceCache = deviceCache
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	err := svcCtx.ServerMsg.Subscribe(eventBus.DmDeviceInfoDelete, func(ctx context.Context, t time.Time, body []byte) error {
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
	logx.Must(err)
	err = svcCtx.ServerMsg.Start()
	logx.Must(err)
}
