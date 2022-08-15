// Package device 设备操作日志
package deviceMsg

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"time"
)

type (
	SDKLog struct {
		ProductID   string    `json:"productID"`  // 产品id
		DeviceName  string    `json:"deviceName"` // 设备名称
		Content     string    `json:"content"`    // 具体信息
		Timestamp   time.Time `json:"timestamp"`  // 操作时间
		ClientToken string    `json:"clientToken"`
		LogLevel    int64     `json:"log_level"`
	}
	SDKLogRepo interface {
		GetDeviceSDKLog(ctx context.Context, productID, deviceName string, page def.PageInfo2) ([]*SDKLog, error)
		Insert(ctx context.Context, data *SDKLog) error
	}
)
