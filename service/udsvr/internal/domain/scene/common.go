package scene

import (
	"context"
	"gitee.com/i-Things/share/caches"
	"github.com/i-Things/things/service/dmsvr/dmExport"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func GetDeviceAlias(ctx context.Context, cache *caches.Cache[dm.DeviceInfo], productID string, deviceName string) string {
	di, err := cache.GetData(ctx, dmExport.GenDeviceInfoKey(productID, deviceName))
	if err != nil {
		return ""
	}
	return di.DeviceAlias.GetValue()
}
