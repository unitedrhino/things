package abnormalLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/stores"
)

func (s *AbnormalLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	return
}

func (s *AbnormalLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&Abnormal{}).Error
	return stores.ErrFmt(err)
}

func (s *AbnormalLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).Where("device_name = ?", deviceName).Delete(&Abnormal{}).Error
	return stores.ErrFmt(err)
}

func (s *AbnormalLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	return nil
}
