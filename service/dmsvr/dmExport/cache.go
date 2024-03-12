package dmExport

import (
	"context"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/eventBus"
	"github.com/i-Things/things/service/dmsvr/client/devicemanage"
	"github.com/i-Things/things/service/dmsvr/client/productmanage"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"strings"
)

func NewProductInfoCache(pm productmanage.ProductManage, fastEvent *eventBus.FastEvent) (*caches.Cache[dm.ProductInfo], error) {
	return caches.NewCache(caches.CacheConfig[dm.ProductInfo]{
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

func NewDeviceInfoCache(devM devicemanage.DeviceManage, fastEvent *eventBus.FastEvent) (*caches.Cache[dm.DeviceInfo], error) {
	return caches.NewCache(caches.CacheConfig[dm.DeviceInfo]{
		KeyType:   eventBus.ServerCacheKeyDmDevice,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*dm.DeviceInfo, error) {
			productID, deviceName, _ := strings.Cut(key, ":")
			ret, err := devM.DeviceInfoRead(ctx, &dm.DeviceInfoReadReq{ProductID: productID, DeviceName: deviceName})
			return ret, err
		},
	})
}
