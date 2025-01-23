package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/domain/schema"
)

func (d *DeviceDataRepo) InitDevice(ctx context.Context,
	t *schema.Model,
	productID string,
	deviceName string) error {
	return nil
}

func (d *DeviceDataRepo) DeleteDevice(
	ctx context.Context,
	t *schema.Model,
	productID string,
	deviceName string) error {
	var typ = map[schema.DataType]struct{}{}
	for _, p := range t.Property {
		if _, ok := typ[p.Define.Type]; ok {
			continue
		}
		typ[p.Define.Type] = struct{}{}
		err := d.db.WithContext(ctx).Table(getTableName(p.Define)).
			Where("product_id = ? and device_name = ?", productID, deviceName).Delete(&Property{}).Error
		if err != nil {
			return stores.ErrFmt(err)
		}
	}
	_, err := d.kv.DelCtx(ctx, d.genRedisPropertyKey(productID, deviceName), d.genRedisPropertyFirstKey(productID, deviceName))
	return err
}
func GetArrayID(id string, num int) string {
	return fmt.Sprintf("%s_%d", id, num)
}

func (d *DeviceDataRepo) DeleteDeviceProperty(ctx context.Context, productID string, deviceName string, s []schema.Property) error {
	var ids []string
	var tables = map[string]struct{}{}
	if len(s) > 0 {
		for _, v := range s {
			ids = append(ids, v.Identifier)
			tables[getTableName(v.Define)] = struct{}{}
		}
		for tb := range tables {
			err := d.db.WithContext(ctx).Table(tb).Where("product_id = ? and device_name = ? and identifier in ?", productID, deviceName, ids).Delete(&Property{}).Error
			if err != nil {
				return stores.ErrFmt(err)
			}
		}

	} else { //删除设备的所有表
		for _, tb := range TableNames {
			err := d.db.WithContext(ctx).Table(tb).Where("product_id = ? and device_name = ?", productID, deviceName).Delete(&Property{}).Error
			if err != nil {
				return stores.ErrFmt(err)
			}
		}
	}
	_, err := d.kv.DelCtx(ctx, d.genRedisPropertyKey(productID, deviceName), d.genRedisPropertyFirstKey(productID, deviceName))
	return err
}
