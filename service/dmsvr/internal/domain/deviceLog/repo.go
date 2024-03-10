// Package device 设备操作日志
package deviceLog

import (
	"context"
	"gitee.com/i-Things/share/devices"
)

type ManageRepo interface {
	InitProduct(ctx context.Context, productID string) error
	InitDevice(ctx context.Context, device devices.Info) error
	DeleteProduct(ctx context.Context, productID string) error
	DeleteDevice(ctx context.Context, productID string, deviceName string) error
}

type ModifyRepo interface {
	ModifyDeviceTenant(ctx context.Context, device devices.Core, tenantCode string) error
	ModifyDeviceArea(ctx context.Context, device devices.Core, areaID int64) error
	ModifyDeviceProject(ctx context.Context, device devices.Core, projectID int64) error
}
