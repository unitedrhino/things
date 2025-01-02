// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"time"
)

type (
	Status struct {
		ProductID  string         `gorm:"column:product_id;type:varchar(100);NOT NULL" json:"productID,omitempty"`   // 产品id
		DeviceName string         `gorm:"column:device_name;type:varchar(100);NOT NULL" json:"deviceName,omitempty"` // 设备名称
		Status     def.ConnStatus `gorm:"column:status;type:BIGINT;NOT NULL" json:"status"`                          // 设备状态 connected:上线 disconnected:下线
		Timestamp  time.Time      `gorm:"column:ts;NOT NULL;" json:"timestamp"`                                      // 操作时间
	}
	StatusFilter struct {
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
