package sendLogRepo

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/devices"
)

func (s *SendLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	s.once.Do(func() {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`user_id` BIGINT,`account` BINARY(200),`action` BINARY(50),`data_id` BINARY(50),`trace_id` BINARY(50),`content` BINARY(200),`result_code` BINARY(50)) "+
			"TAGS (`product_id` BINARY(50),`device_name`  BINARY(50));",
			s.GetLogStableName())
		_, err = s.t.ExecContext(ctx, sql)
	})
	return
}

func (s *SendLogRepo) DeleteProduct(ctx context.Context, productID string) error {
	return nil
}

func (s *SendLogRepo) DeleteDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("drop table if exists %s;", s.GetLogTableName(productID, deviceName))
	if _, err := s.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}

func (s *SendLogRepo) InitDevice(ctx context.Context, device devices.Info) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s (`product_id`,`device_name` ) TAGS (?,?,?,?,?);",
		s.GetLogTableName(device.ProductID, device.DeviceName), s.GetLogStableName())
	_, err := s.t.ExecContext(ctx, sql, device.ProductID, device.DeviceName)
	return err
}
