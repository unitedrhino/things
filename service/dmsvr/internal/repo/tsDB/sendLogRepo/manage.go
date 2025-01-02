package sendLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/stores"
)

func (s *SendLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	return
}

func (s *SendLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&Send{}).Error
	return stores.ErrFmt(err)
}

func (s *SendLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	err := s.db.WithContext(ctx).Where("product_id = ?", productID).
		Where("device_name = ?", deviceName).Delete(&Send{}).Error
	return stores.ErrFmt(err)
}

func (s *SendLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	//sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s (`product_id`,`device_name` ) TAGS (?,?);",
	//	s.GetLogTableName(device.ProductID, device.DeviceName), s.GetLogStableName())
	//_, err := s.t.ExecContext(ctx, sql, device.ProductID, device.DeviceName)
	return nil
}
