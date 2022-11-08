package sdkLogRepo

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/spf13/cast"
)

func ToDeviceSDKLog(productID string, db map[string]any) *deviceMsg.SDKLog {
	return &deviceMsg.SDKLog{
		ProductID:  productID,
		DeviceName: cast.ToString(db["deviceName"]),
		Content:    cast.ToString(db["content"]),
		Timestamp:  cast.ToTime(db["ts"]),
		RequestID:  cast.ToString(db["requestID"]),
		LogLevel:   cast.ToInt64(db["logLevel"]),
	}
}
