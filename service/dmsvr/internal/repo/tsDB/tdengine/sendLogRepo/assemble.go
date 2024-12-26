package sendLogRepo

import (
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceLog(db map[string]any) *deviceLog.Send {
	return &deviceLog.Send{
		UserID:     cast.ToInt64(db["user_id"]),
		Account:    cast.ToString(db["account"]),
		Action:     cast.ToString(db["action"]),
		TraceID:    cast.ToString(db["trace_id"]),
		DataID:     cast.ToString(db["data_id"]),
		Content:    cast.ToString(db["content"]),
		ResultCode: cast.ToInt64(db["result_code"]),
		ProductID:  cast.ToString(db["product_id"]),
		DeviceName: cast.ToString(db["device_name"]),
		Timestamp:  cast.ToTime(db["ts"]),
	}
}
