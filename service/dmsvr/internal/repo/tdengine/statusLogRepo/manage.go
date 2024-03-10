package statusLogRepo

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/devices"
)

func (s *StatusLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	s.once.Do(func() {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`status` BIGINT) "+
			"TAGS (`tenant_code` BINARY(50),`project_id` BIGINT,`area_id` BIGINT,`product_id` BINARY(50),`device_name`  BINARY(50));",
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
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s (`tenant_code`,`project_id`,`area_id`,`product_id`,`device_name` ) TAGS (?,?,?,?,?);",
		s.GetLogTableName(device.ProductID, device.DeviceName), s.GetLogStableName())
	_, err := s.t.ExecContext(ctx, sql, device.TenantCode, device.ProjectID, device.AreaID, device.ProductID, device.DeviceName)
	return err
}

func (s *StatusLogRepo) ModifyDeviceTenant(ctx context.Context, device devices.Core, tenantCode string) error {
	sql := fmt.Sprintf("ALTER TABLE %s SET TAG `tenant_code`=?;",
		s.GetLogTableName(device.ProductID, device.DeviceName))
	_, err := s.t.ExecContext(ctx, sql, tenantCode)
	return err
}

func (s *StatusLogRepo) ModifyDeviceArea(ctx context.Context, device devices.Core, areaID int64) error {
	sql := fmt.Sprintf("ALTER TABLE %s SET TAG `area_id`=?;",
		s.GetLogTableName(device.ProductID, device.DeviceName))
	_, err := s.t.ExecContext(ctx, sql, areaID)
	return err
}

func (s *StatusLogRepo) ModifyDeviceProject(ctx context.Context, device devices.Core, projectID int64) error {
	sql := fmt.Sprintf("ALTER TABLE %s SET TAG `project_id`=?;",
		s.GetLogTableName(device.ProductID, device.DeviceName))
	_, err := s.t.ExecContext(ctx, sql, projectID)
	return err
}
