package statusLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/devices"
)

func (s *StatusLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	return
}

func (s *StatusLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&Status{}).Error
	return stores.ErrFmt(err)
}

func (s *StatusLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).
		Where("device_name = ?", deviceName).Delete(&Status{}).Error
	return stores.ErrFmt(err)
}

func (s *StatusLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	return nil
}
