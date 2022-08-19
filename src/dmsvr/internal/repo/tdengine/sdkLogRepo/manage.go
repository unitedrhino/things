package sdkLogRepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d SDKLogRepo) InitProduct(ctx context.Context, productID string) error {
	sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
		"(`ts` timestamp,`content` BINARY(5000),`log_level` BINARY(100),`client_token` BINARY(100))"+
		"TAGS (device_name BINARY(50));",
		d.GetSDKLogStableName(productID))
	if _, err := d.t.ExecContext(ctx,sql); err != nil {
		return err
	}
	return nil
}

func (d SDKLogRepo) InitDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS ('%s');",
		d.GetSDKLogTableName(productID, deviceName), d.GetSDKLogStableName(productID), deviceName)
	if _, err := d.t.ExecContext(ctx,sql); err != nil {
		logx.WithContext(ctx).Errorf(
			"%s|EventTable|productID:%v,deviceName:%v,err:%v",
			utils.FuncName(), productID, deviceName, err)
		return err
	}
	return nil
}

func (d SDKLogRepo) DropProduct(ctx context.Context, productID string) error {
	sql := fmt.Sprintf("drop stable if exists %s;", d.GetSDKLogStableName(productID))
	if _, err := d.t.ExecContext(ctx,sql); err != nil {
		return err
	}
	return nil
}

func (d SDKLogRepo) DropDevice(ctx context.Context, productID string, deviceName string) error {
	sql := fmt.Sprintf("drop table if exists %s;", d.GetSDKLogTableName(productID, deviceName))
	if _, err := d.t.ExecContext(ctx,sql); err != nil {
		return err
	}
	return nil
}
