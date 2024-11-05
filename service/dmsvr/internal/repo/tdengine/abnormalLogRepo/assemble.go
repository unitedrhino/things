package abnormalLogRepo

import (
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceLog(db map[string]any) *deviceLog.Abnormal {
	return &deviceLog.Abnormal{
		Type:       cast.ToString(db["type"]),
		Reason:     cast.ToString(db["reason"]),
		Action:     cast.ToBool(db["action"]),
		TraceID:    cast.ToString(db["trace_id"]),
		ProductID:  cast.ToString(db["product_id"]),
		DeviceName: cast.ToString(db["device_name"]),
		Timestamp:  cast.ToTime(db["ts"]),
	}
}
