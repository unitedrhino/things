package sendLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"gitee.com/unitedrhino/things/share/devices"
)

func (s *SendLogRepo) InitProduct(ctx context.Context, productID string) (err error) {
	s.once.Do(func() {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`user_id` BIGINT,`account` BINARY(200),`action` BINARY(50),`data_id` BINARY(50),`trace_id` BINARY(50),`content` BINARY(200),`result_code` BINARY(50)) "+
			"TAGS (%s);",
			s.GetLogStableName(), tdengine.GenTagsDef(defaultTagDef, s.groupConfigs))
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
	//sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s (`product_id`,`device_name` ) TAGS (?,?);",
	//	s.GetLogTableName(device.ProductID, device.DeviceName), s.GetLogStableName())
	//_, err := s.t.ExecContext(ctx, sql, device.ProductID, device.DeviceName)
	return nil
}

func (s *SendLogRepo) UpdateDevice(ctx context.Context, devices []*devices.Core, affiliation devices.Affiliation) error {
	var tables []string
	for _, device := range devices {
		tables = append(tables, s.GetLogTableName(device.ProductID, device.DeviceName))
	}
	err := tdengine.AlterTag(ctx, s.t, tables, tdengine.AffiliationToMap(affiliation, nil))
	return err
}

func (s *SendLogRepo) VersionUpdate(ctx context.Context, version string) error {
	s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s DROP TAG `group_ids` ;", s.GetLogStableName()))
	s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s DROP TAG `group_id_paths`;", s.GetLogStableName()))
	s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s ADD TAG `tenant_code`  BINARY(50) ;", s.GetLogStableName()))
	s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s ADD TAG  `project_id` BIGINT ;", s.GetLogStableName()))
	s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s ADD TAG  `area_id` BIGINT  ;", s.GetLogStableName()))
	s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s ADD TAG `area_id_path`  BINARY(50) ;", s.GetLogStableName()))
	for _, g := range s.groupConfigs {
		s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s ADD TAG `group_%s_ids`  BINARY(250) ;", s.GetLogStableName(), g.Value))
		s.t.ExecContext(ctx, fmt.Sprintf("ALTER STABLE %s ADD TAG `group_%s_id_paths`  BINARY(250) ;", s.GetLogStableName(), g.Value))
	}
	return nil
}

func (s *SendLogRepo) UpdateDevices(ctx context.Context, devs []*devices.Info) error {
	var tags []tdengine.Tag
	for _, device := range devs {
		tags = append(tags, tdengine.Tag{
			Table: s.GetLogTableName(device.ProductID, device.DeviceName),
			Tags:  tdengine.AffiliationToMap(utils.Copy2[devices.Affiliation](device), nil),
		})
	}
	err := tdengine.AlterTags(ctx, s.t, tags)
	return err
}
