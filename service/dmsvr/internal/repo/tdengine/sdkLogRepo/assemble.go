package sdkLogRepo

import (
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceSDKLog(productID string, db map[string]any) *deviceLog.SDK {
	return &deviceLog.SDK{
		ProductID:  productID,
		DeviceName: cast.ToString(db["device_name"]),
		Content:    cast.ToString(db["content"]),
		Timestamp:  cast.ToTime(db["ts"]),
		LogLevel:   cast.ToInt64(db["log_level"]),
	}
}
