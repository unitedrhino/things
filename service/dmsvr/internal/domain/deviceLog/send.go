// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/i-Things/share/def"
	"time"
)

type (
	Send struct {
		TenantCode string    `json:"tenantCode,omitempty"`
		ProjectID  int64     `json:"projectID,omitempty"`
		AreaID     int64     `json:"areaID"`
		UserID     int64     `json:"userID"`
		ProductID  string    `json:"productID,omitempty"`  // 产品id
		DeviceName string    `json:"deviceName,omitempty"` // 设备名称
		Action     string    `json:"action,omitempty"`     // 操作类型 propertySend:属性控制 actionSend:操作控制 propertyGetReportSend:获取最新属性请求
		DataID     string    `json:"dataID"`
		Timestamp  time.Time `json:"timestamp"`            // 操作时间
		TraceID    string    `json:"traceID,omitempty"`    // 服务器端事务id
		Content    string    `json:"content"`              //操作的内容
		ResultCode int64     `json:"resultCode,omitempty"` // 请求结果状态,200为成功
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
		ModifyRepo
	}
)
