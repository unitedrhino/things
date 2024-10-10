// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"time"
)

type (
	Status struct {
		ProductID  string    `json:"productID"`  // 产品id
		DeviceName string    `json:"deviceName"` // 设备名称
		Status     int64     `json:"status"`     // 设备状态 connected:上线 disconnected:下线
		Timestamp  time.Time `json:"timestamp"`  // 操作时间
	}
	StatusFilter struct {
		TenantCode string
		ProjectID  int64
		AreaIDs    []int64
		ProductID  string // 产品id
		DeviceName string // 设备名称
		Status     int64  `json:"status"`
	}
	StatusRepo interface {
		GetDeviceLog(ctx context.Context, filter StatusFilter, page def.PageInfo2) ([]*Status, error)
		GetCountLog(ctx context.Context, filter StatusFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Status) error
		ManageRepo
		//ModifyRepo
	}
)
