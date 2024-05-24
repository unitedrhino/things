package scene

import (
	"context"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/devices"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func GetDeviceAlias(ctx context.Context, cache *caches.Cache[dm.DeviceInfo, devices.Core], productID string, deviceName string) string {
	di, err := cache.GetData(ctx, devices.Core{
		ProductID:  productID,
		DeviceName: deviceName,
	})
	if err != nil {
		return ""
	}
	return di.DeviceAlias.GetValue()
}
