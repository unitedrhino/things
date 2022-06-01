// Package device 设备操作日志
package device

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
		TranceID    string    `json:"tranceID"`   // 服务器端事务id
		ResultType  int64     `json:"resultType"` // 请求结果状态,0为成功
	}
	SDKLogRepo interface {
		GetDeviceDebugLog(ctx context.Context, productID, deviceName string, page def.PageInfo2) ([]*SDKLog, error)
		Insert(ctx context.Context, data *SDKLog) error
		InitProduct(ctx context.Context, productID string) error
		InitDevice(ctx context.Context, productID string, deviceName string) error
		DropProduct(ctx context.Context, productID string) error
		DropDevice(ctx context.Context, productID string, deviceName string) error
	}
)
