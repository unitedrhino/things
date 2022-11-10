package sdkLogRepo

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgSdkLog"
	"github.com/spf13/cast"
)

func ToDeviceSDKLog(productID string, db map[string]any) *msgSdkLog.SDKLog {
	return &msgSdkLog.SDKLog{
		ProductID:  productID,
		DeviceName: cast.ToString(db["deviceName"]),
		Content:    cast.ToString(db["content"]),
		Timestamp:  cast.ToTime(db["ts"]),
		LogLevel:   cast.ToInt64(db["logLevel"]),
	}
}
