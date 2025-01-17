package abnormalLogRepo

import (
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"github.com/spf13/cast"
)

func ToDeviceLog(db map[string]any) *deviceLog.Abnormal {
	return &deviceLog.Abnormal{
		TenantCode: stores.TenantCode(cast.ToString(db["tenant_code"])),
		ProjectID:  stores.ProjectID(cast.ToInt64(db["project_id"])),
		AreaID:     stores.AreaID(cast.ToInt64(db["area_id"])),
		AreaIDPath: stores.AreaIDPath(cast.ToString(db["area_id_path"])),
		Type:       cast.ToString(db["type"]),
		Reason:     cast.ToString(db["reason"]),
		Action:     def.ToIntBool[int64](cast.ToBool(db["action"])),
		TraceID:    cast.ToString(db["trace_id"]),
		ProductID:  cast.ToString(db["product_id"]),
		DeviceName: cast.ToString(db["device_name"]),
		Timestamp:  cast.ToTime(db["ts"]),
	}
}
