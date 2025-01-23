package sdkLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/devices"
)

func (d SDKLogRepo) InitProduct(ctx context.Context, productID string) error {
	return nil
}

func (d SDKLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	return nil
}

func (d SDKLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	err := d.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&SDK{}).Error
	return stores.ErrFmt(err)
}

func (d SDKLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	err := d.db.WithContext(ctx).Where("product_id = ?", productID).
		Where("device_name = ?", deviceName).Delete(&SDK{}).Error
	return stores.ErrFmt(err)
}
