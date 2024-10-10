package hubLogRepo

import (
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceLog(productID string, db map[string]any) *deviceLog.Hub {
	return &deviceLog.Hub{
		ProductID:   productID,
		DeviceName:  cast.ToString(db["device_name"]),
		Content:     cast.ToString(db["content"]),
		Topic:       cast.ToString(db["topic"]),
		Action:      cast.ToString(db["action"]),
		Timestamp:   cast.ToTime(db["ts"]),
		RequestID:   cast.ToString(db["request_id"]),
		TraceID:     cast.ToString(db["trace_id"]),
		ResultCode:  cast.ToInt64(db["result_type"]),
		RespPayload: cast.ToString(db["resp_payload"]),
	}
}
