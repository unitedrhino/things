package statusLogRepo

import (
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceLog(db map[string]any) *deviceLog.Status {
	return &deviceLog.Status{
		Status:     cast.ToInt64(db["status"]),
		ProductID:  cast.ToString(db["product_id"]),
		DeviceName: cast.ToString(db["device_name"]),
		Timestamp:  cast.ToTime(db["ts"]),
	}
}
