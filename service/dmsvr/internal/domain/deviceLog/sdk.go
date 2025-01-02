// Package device 设备操作日志
package deviceLog

import (
	"context"
	"time"

	"gitee.com/unitedrhino/share/def"
)

type (
	SDK struct {
		ProductID  string    `gorm:"column:product_id;type:varchar(100);NOT NULL" json:"productID,omitempty"`   // 产品id
		DeviceName string    `gorm:"column:device_name;type:varchar(100);NOT NULL" json:"deviceName,omitempty"` // 设备名称
		Content    string    `gorm:"column:content;type:varchar(256);NOT NULL" json:"content,omitempty"`        // 具体信息
		Timestamp  time.Time `gorm:"column:ts;NOT NULL;" json:"timestamp"`                                      // 操作时间
		LogLevel   int64     `gorm:"column:log_level;type:BIGINT;default:1" json:"logLevel"`
	}
	SDKFilter struct {
		ProductID  string // 产品id
		DeviceName string // 设备名称
		LogLevel   int    //日志等级
	}
	SDKRepo interface {
		GetDeviceSDKLog(ctx context.Context, filter SDKFilter, page def.PageInfo2) ([]*SDK, error)
		GetCountLog(ctx context.Context, filter SDKFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *SDK) error
		ManageRepo
	}
)
