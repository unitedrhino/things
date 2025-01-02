// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"time"
)

type (
	Send struct {
		UserID     int64     `gorm:"column:user_id;type:BIGINT;NOT NULL" json:"userID"`
		ProductID  string    `gorm:"column:product_id;type:varchar(100);NOT NULL" json:"productID,omitempty"`   // 产品id
		DeviceName string    `gorm:"column:device_name;type:varchar(100);NOT NULL" json:"deviceName,omitempty"` // 设备名称
		Action     string    `gorm:"column:action;type:varchar(100);NOT NULL" json:"action,omitempty"`          // 操作类型 propertySend:属性控制 actionSend:操作控制 propertyGetReportSend:获取最新属性请求
		DataID     string    `gorm:"column:data_id;type:varchar(100);NOT NULL" json:"dataID"`
		Timestamp  time.Time `gorm:"column:ts;NOT NULL;" json:"timestamp"`                                // 操作时间
		TraceID    string    `gorm:"column:trace_id;type:varchar(100);NOT NULL" json:"traceID,omitempty"` // 服务器端事务id
		Account    string    `gorm:"column:account;type:varchar(100);NOT NULL" json:"account"`
		Content    string    `gorm:"column:content;type:varchar(100);NOT NULL" json:"content"`               //操作的内容
		ResultCode int64     `gorm:"column:result_code;type:BIGINT;default:200" json:"resultCode,omitempty"` // 请求结果状态
	}
	SendFilter struct {
		TenantCode string
		ProjectID  int64    `json:"projectID,omitempty"`
		AreaIDs    []int64  `json:"areaID"`
		UserID     int64    `json:"userID"`
		ProductID  string   // 产品id
		DeviceName string   // 设备名称
		Actions    []string //过滤操作类型  propertySend:属性控制 actionSend:操作控制 propertyGetReportSend:获取最新属性请求
		ResultCode int64
	}

	SendRepo interface {
		GetDeviceLog(ctx context.Context, filter SendFilter, page def.PageInfo2) ([]*Send, error)
		GetCountLog(ctx context.Context, filter SendFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Send) error
		ManageRepo
		//ModifyRepo
	}
)
