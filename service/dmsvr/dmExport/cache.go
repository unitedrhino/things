package dmExport

import (
	"context"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/eventBus"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func NewProductInfoCache(pm productmanage.ProductManage, fastEvent *eventBus.FastEvent) (*caches.Cache[dm.ProductInfo, string], error) {
	return caches.NewCache(caches.CacheConfig[dm.ProductInfo, string]{
		KeyType:   eventBus.ServerCacheKeyDmProduct,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*dm.ProductInfo, error) {
			ret, err := pm.ProductInfoRead(ctx, &dm.ProductInfoReadReq{ProductID: key, WithCategory: true, WithProtocol: true})
			return ret, err
		},
	})
}

func GenDeviceInfoKey(productID, deviceName string) string {
	return productID + ":" + deviceName
}

func NewDeviceInfoCache(devM devicemanage.DeviceManage, fastEvent *eventBus.FastEvent) (*caches.Cache[dm.DeviceInfo, devices.Core], error) {
	return caches.NewCache(caches.CacheConfig[dm.DeviceInfo, devices.Core]{
		KeyType:   eventBus.ServerCacheKeyDmDevice,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key devices.Core) (*dm.DeviceInfo, error) {
			ret, err := devM.DeviceInfoRead(ctx, &dm.DeviceInfoReadReq{ProductID: key.ProductID, DeviceName: key.DeviceName})
			return ret, err
		},
	})
}

func NewSchemaInfoCache(pm productmanage.ProductManage, fastEvent *eventBus.FastEvent) (*caches.Cache[schema.Model, string], error) {
	return caches.NewCache(caches.CacheConfig[schema.Model, string]{
		KeyType:   eventBus.ServerCacheKeyDmSchema,
		FastEvent: fastEvent,
		Fmt: func(ctx context.Context, key string, data *schema.Model) {
			data.ValidateWithFmt()
		},
		GetData: func(ctx context.Context, key string) (*schema.Model, error) {
			info, err := pm.ProductSchemaTslRead(ctx, &dm.ProductSchemaTslReadReq{ProductID: key})
			if err != nil {
				return nil, err
			}
			return schema.ValidateWithFmt([]byte(info.Tsl))
		},
	})
}
