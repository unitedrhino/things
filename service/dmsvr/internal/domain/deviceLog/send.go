// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/things/share/devices"
	"time"
)

type (
	Send struct {
		TenantCode dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                        // 租户编码
		ProjectID  dataType.ProjectID  `gorm:"column:project_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
		AreaID     dataType.AreaID     `gorm:"column:area_id;index:project_id_area_id;type:bigint;default:0;NOT NULL"`    // 项目区域ID(雪花ID)
		AreaIDPath dataType.AreaIDPath `gorm:"column:area_id_path;type:varchar(100);default:'';NOT NULL"`                 // 项目区域ID路径(雪花ID)
		UserID     int64               `gorm:"column:user_id;type:BIGINT;NOT NULL" json:"userID"`
		ProductID  string              `gorm:"column:product_id;type:varchar(100);NOT NULL" json:"productID,omitempty"`   // 产品id
		DeviceName string              `gorm:"column:device_name;type:varchar(100);NOT NULL" json:"deviceName,omitempty"` // 设备名称
		Action     string              `gorm:"column:action;type:varchar(100);NOT NULL" json:"action,omitempty"`          // 操作类型 propertySend:属性控制 actionSend:操作控制 propertyGetReportSend:获取最新属性请求
		DataID     string              `gorm:"column:data_id;type:varchar(100);NOT NULL" json:"dataID"`
		Timestamp  time.Time           `gorm:"column:ts;NOT NULL;" json:"timestamp"`                                // 操作时间
		TraceID    string              `gorm:"column:trace_id;type:varchar(100);NOT NULL" json:"traceID,omitempty"` // 服务器端事务id
		Account    string              `gorm:"column:account;type:varchar(100);NOT NULL" json:"account"`
		Content    string              `gorm:"column:content;type:varchar(100);NOT NULL" json:"content"`               //操作的内容
		ResultCode int64               `gorm:"column:result_code;type:BIGINT;default:200" json:"resultCode,omitempty"` // 请求结果状态
	}
	SendFilter struct {
		TenantCode string
		ProjectID  int64   `json:"projectID,omitempty"`
		AreaID     int64   `json:"areaID,omitempty"`
		AreaIDPath string  `json:"areaIDPath,omitempty"`
		AreaIDs    []int64 `json:"areaIDs"`
		UserID     int64   `json:"userID"`
		ProductIDs []string
		ProductID  string   // 产品id
		DeviceName string   // 设备名称
		Actions    []string //过滤操作类型  propertySend:属性控制 actionSend:操作控制 propertyGetReportSend:获取最新属性请求
		DataIDs    []string
		DataID     string
		ResultCode int64
	}

	SendRepo interface {
		GetDeviceLog(ctx context.Context, filter SendFilter, page def.PageInfo2) ([]*Send, error)
		GetCountLog(ctx context.Context, filter SendFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Send) error
		ManageRepo
		UpdateDevice(ctx context.Context, devices []*devices.Core, affiliation devices.Affiliation) error
		UpdateDevices(ctx context.Context, devices []*devices.Info) error
		VersionUpdate(ctx context.Context, version string) error
		//ModifyRepo
	}
)
