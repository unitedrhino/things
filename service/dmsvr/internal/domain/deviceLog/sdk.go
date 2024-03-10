// Package device 设备操作日志
package deviceLog

import (
	"context"
	"time"

	"gitee.com/i-Things/share/def"
)

type (
	SDK struct {
		ProductID  string    `json:"productID"`  // 产品id
		DeviceName string    `json:"deviceName"` // 设备名称
		Content    string    `json:"content"`    // 具体信息
		Timestamp  time.Time `json:"timestamp"`  // 操作时间
		LogLevel   int64     `json:"logLevel"`
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
