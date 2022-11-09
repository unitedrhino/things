// Package device 设备操作日志
package msgSdkLog

import (
	"context"
	"time"

	"github.com/i-Things/things/shared/def"
)

type (
	SDKLog struct {
		ProductID  string    `json:"productID"`  // 产品id
		DeviceName string    `json:"deviceName"` // 设备名称
		Content    string    `json:"content"`    // 具体信息
		Timestamp  time.Time `json:"timestamp"`  // 操作时间
		LogLevel   int64     `json:"logLevel"`
	}
	SdkLogFilter struct {
		ProductID  string // 产品id
		DeviceName string // 设备名称
		LogLevel   int    //日志等级
	}
	SDKLogRepo interface {
		GetDeviceSDKLog(ctx context.Context, filter SdkLogFilter, page def.PageInfo2) ([]*SDKLog, error)
		GetCountLog(ctx context.Context, filter SdkLogFilter, page def.PageInfo2) (int64, error)
		Insert(ctx context.Context, data *SDKLog) error
	}
)
