// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"time"
)

type (
	Abnormal struct {
		ProductID  string    `gorm:"column:product_id;type:varchar(100);NOT NULL" json:"productID,omitempty"`   // 产品id
		DeviceName string    `gorm:"column:device_name;type:varchar(100);NOT NULL" json:"deviceName,omitempty"` // 设备名称
		Action     def.Bool  `gorm:"column:user_id;type:BIGINT;default:1" json:"action,omitempty"`              //触发1还是解除2
		Type       string    `gorm:"column:type;type:varchar(100);NOT NULL" json:"type,omitempty"`              // 异常类型
		Timestamp  time.Time `gorm:"column:ts;NOT NULL;" json:"timestamp"`                                      // 操作时间
		TraceID    string    `gorm:"column:trace_id;type:varchar(100);NOT NULL" json:"traceID,omitempty"`       // 服务器端事务id
		Reason     string    `gorm:"column:reason;type:varchar(100);NOT NULL" json:"reason,omitempty"`          //原因
	}
	AbnormalFilter struct {
		ProductID  string // 产品id
		DeviceName string // 设备名称
		Action     int64
		Type       string `json:"type,omitempty"` // 异常类型
		Reason     string
	}

	AbnormalRepo interface {
		GetDeviceLog(ctx context.Context, filter AbnormalFilter, page def.PageInfo2) ([]*Abnormal, error)
		GetCountLog(ctx context.Context, filter AbnormalFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *Abnormal) error
		ManageRepo
		//ModifyRepo
	}
)
