// Package device 设备操作日志
package deviceMsgManage

import (
	"context"
)

type (
	SDKLogRepo interface {
		InitProduct(ctx context.Context, productID string) error
		InitDevice(ctx context.Context, productID string, deviceName string) error
		DropProduct(ctx context.Context, productID string) error
		DropDevice(ctx context.Context, productID string, deviceName string) error
	}
)
