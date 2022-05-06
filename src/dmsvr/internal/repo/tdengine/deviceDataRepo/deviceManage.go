package deviceDataRepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/templateModel"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) InitDevice(ctx context.Context,
	t *templateModel.Template,
	productID string,
	deviceName string) error {
	err := d.createPropertyTable(ctx, t.Properties, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf(
			"%s|createPropertyTable|productID:%v,deviceName:%v,err:%v,properties:%v",
			utils.FuncName(), productID, deviceName, t.Properties, err)
		return err
	}
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS ('%s');",
		getEventTableName(productID, deviceName), getEventStableName(productID), deviceName)
	if _, err := d.t.Exec(sql); err != nil {
		logx.WithContext(ctx).Errorf(
			"%s|EventTable|productID:%v,deviceName:%v,err:%v",
			utils.FuncName(), productID, deviceName, err)
		return err
	}
	return nil
}

func (d *DeviceDataRepo) DropDevice(
	ctx context.Context,
	t *templateModel.Template,
	productID string,
	deviceName string) error {
	tableList := getTableNameList(t, productID, deviceName)
	for _, v := range tableList {
		sql := fmt.Sprintf("drop table if exists %s;", v)
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeviceDataRepo) createPropertyTable(
	ctx context.Context, p templateModel.Properties, productID string, deviceName string) error {
	for _, v := range p {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s  TAGS('%s','%s');",
			getPropertyTableName(productID, deviceName, v.ID),
			getPropertyStableName(productID, v.ID), deviceName, v.Define.Type)
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
