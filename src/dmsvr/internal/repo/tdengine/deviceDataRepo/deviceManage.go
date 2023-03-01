package deviceDataRepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) InitDevice(ctx context.Context,
	t *schema.Model,
	productID string,
	deviceName string) error {
	if t.Property != nil { //如果还没有定义属性则不需要初始化
		err := d.createPropertyTable(ctx, t.Property, productID, deviceName)
		if err != nil {
			logx.WithContext(ctx).Errorf(
				"%s.createPropertyTable productID:%v,deviceName:%v,err:%v,properties:%v",
				utils.FuncName(), productID, deviceName, t.Property, err)
			return err
		}
	}
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS ('%s','%s');",
		d.GetEventTableName(productID, deviceName), d.GetEventStableName(), productID, deviceName)
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		logx.WithContext(ctx).Errorf(
			"%s.EventTable productID:%v,deviceName:%v,err:%v",
			utils.FuncName(), productID, deviceName, err)
		return err
	}
	return nil
}

func (d *DeviceDataRepo) DeleteDevice(
	ctx context.Context,
	t *schema.Model,
	productID string,
	deviceName string) error {
	tableList := d.GetTableNameList(t, productID, deviceName)
	for _, v := range tableList {
		sql := fmt.Sprintf("drop table if exists %s;", v)
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeviceDataRepo) createPropertyTable(
	ctx context.Context, p schema.PropertyMap, productID string, deviceName string) error {
	for _, v := range p {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s  TAGS('%s','%s');",
			d.GetPropertyTableName(productID, deviceName, v.Identifier),
			d.GetPropertyStableName(productID, v.Identifier), deviceName, v.Define.Type)
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return err
		}
	}
	return nil
}
