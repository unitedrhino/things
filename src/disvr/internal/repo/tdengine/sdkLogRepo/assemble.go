package sdkLogRepo

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/spf13/cast"
)

func ToDeviceSDKLog(productID string, db map[string]any) *deviceMsg.SDKLog {
	return &deviceMsg.SDKLog{
		ProductID:   productID,
		DeviceName:  cast.ToString(db["device_name"]),
		Content:     cast.ToString(db["content"]),
		Timestamp:   cast.ToTime(db["ts"]),
		ClientToken: cast.ToString(db["client_token"]),
		LogLevel:    cast.ToInt64(db["log_level"]),
	}
}
