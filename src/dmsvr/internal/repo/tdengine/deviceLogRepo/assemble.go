package deviceLogRepo

import (
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/spf13/cast"
)

func ToDeviceLog(productID string, db map[string]interface{}) *device.Log {
	return &device.Log{
		ProductID:  productID,
		DeviceName: cast.ToString(db["device_name"]),
		Content:    cast.ToString(db["content"]),
		Topic:      cast.ToString(db["topic"]),
		Action:     cast.ToString(db["action"]),
		Timestamp:  cast.ToTime(db["ts"]),
		RequestID:  cast.ToString(db["request_id"]),
		TranceID:   cast.ToString(db["trance_id"]),
		ResultType: cast.ToInt64(db["result_type"]),
	}
}
