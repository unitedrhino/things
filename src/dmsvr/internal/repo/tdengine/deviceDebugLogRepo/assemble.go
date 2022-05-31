package deviceDebugLogRepo

import (
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/spf13/cast"
)

func ToDeviceDebugLog(productID string, db map[string]interface{}) *device.DebugLog {
	return &device.DebugLog{
		ProductID:  productID,
		DeviceName: cast.ToString(db["device_name"]),
		Content:    cast.ToString(db["content"]),
		Timestamp:  cast.ToTime(db["ts"]),
		RequestID:  cast.ToString(db["request_id"]),
		TranceID:   cast.ToString(db["trance_id"]),
		ResultType: cast.ToInt64(db["result_type"]),
	}
}
