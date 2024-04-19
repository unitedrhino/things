package sendLogRepo

import (
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceLog(db map[string]any) *deviceLog.Send {
	return &deviceLog.Send{
		TenantCode: cast.ToString(db["tenant_code"]),
		ProjectID:  cast.ToInt64(db["project_id"]),
		AreaID:     cast.ToInt64(db["area_id"]),
		UserID:     cast.ToInt64(db["user_id"]),
		Action:     cast.ToString(db["action"]),
		TraceID:    cast.ToString(db["trace_id"]),
		ResultCode: cast.ToInt64(db["result_code"]),
		ProductID:  cast.ToString(db["product_id"]),
		DeviceName: cast.ToString(db["device_name"]),
		Timestamp:  cast.ToTime(db["ts"]),
	}
}
