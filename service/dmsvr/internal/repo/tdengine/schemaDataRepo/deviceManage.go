package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func (d *DeviceDataRepo) InitDevice(ctx context.Context,
	t *schema.Model,
	productID string,
	deviceName string) error {
	err := d.createPropertyTable(ctx, t.Property, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Errorf(
			"%s.createPropertyTable productID:%v,deviceName:%v,err:%v,properties:%v",
			utils.FuncName(), productID, deviceName, utils.Fmt(t.Property), err)
		return err
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
	var sqls []string
	for _, v := range tableList {
		sqls = append(sqls, fmt.Sprintf(" if exists %s ", v))
	}
	sql := fmt.Sprintf("DROP TABLE %s", strings.Join(sqls, ","))
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	err := d.DeleteDeviceProperty(ctx, productID, deviceName, nil)
	if err != nil {
		return err
	}
	_, err = d.kv.DelCtx(ctx, d.genRedisPropertyKey(productID, deviceName), d.genRedisPropertyFirstKey(productID, deviceName))
	return err
}
func GetArrayID(id string, num int) string {
	return fmt.Sprintf("%s_%d", id, num)
}

func (d *DeviceDataRepo) createPropertyTable(
	ctx context.Context, p schema.PropertyMap, productID string, deviceName string) error {
	for _, v := range p {
		if v.Define.Type == schema.DataTypeArray {
			for i := 0; i < cast.ToInt(v.Define.Max); i++ {
				sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s  TAGS('%s','%s',%d,'%s');",
					d.GetPropertyTableName(productID, deviceName, GetArrayID(v.Identifier, i)),
					d.GetPropertyStableName(v.Tag, productID, deviceName, v.Identifier), productID, deviceName, i, v.Define.ArrayInfo.Type)
				if _, err := d.t.ExecContext(ctx, sql); err != nil {
					return err
				}
			}
		} else {
			sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s USING %s  TAGS('%s','%s','%s');",
				d.GetPropertyTableName(productID, deviceName, v.Identifier),
				d.GetPropertyStableName(v.Tag, productID, deviceName, v.Identifier), productID, deviceName, v.Define.Type)
			if _, err := d.t.ExecContext(ctx, sql); err != nil {
				return err
			}
		}

	}
	return nil
}

func (d *DeviceDataRepo) DeleteDeviceProperty(ctx context.Context, productID string, deviceName string, s []schema.Property) error {
	var sqls []string
	if len(s) > 0 {
		for _, v := range s {
			sqls = append(sqls, fmt.Sprintf(" if exists %s ", d.GetPropertyTableName(productID, deviceName, v.Identifier)))
		}
	} else { //删除设备的所有表
		rows, err := d.t.QueryContext(ctx, fmt.Sprintf("SHOW TABLE TAGS FROM %s;", d.GetDeviceStableFloatName()))
		if err != nil {
			return err
		}
		var datas []map[string]any
		stores.Scan(rows, &datas)
	}
	sql := fmt.Sprintf("DROP TABLE %s", strings.Join(sqls, ","))
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return err
	}
	_, err := d.kv.DelCtx(ctx, d.genRedisPropertyKey(productID, deviceName), d.genRedisPropertyFirstKey(productID, deviceName))
	return err
}
