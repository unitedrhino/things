package hubLogRepo

import (
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/spf13/cast"
)

func ToDeviceLog(productID string, db map[string]any) *deviceMsg.HubLog {
	return &deviceMsg.HubLog{
		ProductID:  productID,
		DeviceName: cast.ToString(db["deviceName"]),
		Content:    cast.ToString(db["content"]),
		Topic:      cast.ToString(db["topic"]),
		Action:     cast.ToString(db["action"]),
		Timestamp:  cast.ToTime(db["ts"]),
		RequestID:  cast.ToString(db["requestID"]),
		TranceID:   cast.ToString(db["trance_id"]),
		ResultType: cast.ToInt64(db["result_type"]),
	}
}
