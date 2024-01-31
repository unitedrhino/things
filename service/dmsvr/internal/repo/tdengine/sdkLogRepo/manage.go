package sdkLogRepo

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d SDKLogRepo) InitProduct(ctx context.Context, productID string) error {
	sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
		"(`ts` timestamp,`content` BINARY(5000),`log_level` BINARY(100)) "+
		" TAGS (`product_id` BINARY(50), `device_name` BINARY(50));",
		d.GetSDKLogStableName())
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}

func (d SDKLogRepo) InitDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS ('%s','%s');",
		d.GetSDKLogTableName(productID, deviceName), d.GetSDKLogStableName(), productID, deviceName)
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		logx.WithContext(ctx).Errorf(
			"%s.ExecContext productID:%v,deviceName:%v,err:%v",
			utils.FuncName(), productID, deviceName, err)
		return err
	}
	return nil
}

func (d SDKLogRepo) DropProduct(ctx context.Context, productID string) error {
	return nil
}

func (d SDKLogRepo) DropDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("drop table if exists %s;", d.GetSDKLogTableName(productID, deviceName))
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}
