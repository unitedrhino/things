package statusLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/devices"
)

func (s *StatusLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	s.once.Do(func() {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`status` BIGINT) "+
			"TAGS (`product_id` BINARY(50),`device_name`  BINARY(50));",
			s.GetLogStableName())
		_, err = s.t.ExecContext(ctx, sql)
	})
	return
}

func (s *StatusLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	return nil
}

func (s *StatusLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("drop table if exists %s;", s.GetLogTableName(productID, deviceName))
	if _, err := s.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}

func (s *StatusLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s (`product_id`,`device_name` ) TAGS (?,?);",
		s.GetLogTableName(device.ProductID, device.DeviceName), s.GetLogStableName())
	_, err := s.t.ExecContext(ctx, sql, device.ProductID, device.DeviceName)
	return err
}
