package tdengine

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) InitDevice(ctx context.Context, t *deviceTemplate.Template, productID string, deviceName string) error {
	err := d.createPropertyTable(ctx, t.Properties, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s|createPropertyTable|productID:%v,deviceName:%v,properties:%v,err:%v",
			utils.FuncName(), productID, deviceName, t.Properties, err)
		return err
	}
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS ('%s');",
		getEventTableName(productID, deviceName), getEventStableName(productID), deviceName)
	if _, err := d.t.Exec(sql); err != nil {
		return err
	}
	sql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS ('%s');",
		getActionTableName(productID, deviceName), getActionStableName(productID), deviceName)
	if _, err := d.t.Exec(sql); err != nil {
		return err
	}
	return nil
}

func (d *DeviceDataRepo) DropDevice(ctx context.Context, t *deviceTemplate.Template, productID string, deviceName string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) createPropertyTable(
	ctx context.Context, p deviceTemplate.Properties, productID string, deviceName string) error {
	for _, v := range p {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s TAGS('%s');",
			getPropertyTableName(productID, deviceName, v.ID), getPropertyStableName(productID, v.ID), deviceName)
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
