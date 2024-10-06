package scene

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/things/service/dmsvr/dmExport"
)

func GetDeviceAlias(ctx context.Context, cache dmExport.DeviceCacheT, productID string, deviceName string) string {
	di, err := cache.GetData(ctx, devices.Core{
		ProductID:  productID,
		DeviceName: deviceName,
	})
	if err != nil {
		return ""
	}
	return di.DeviceAlias.GetValue()
}
