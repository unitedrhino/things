package schemaDataRepo

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
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
	err := d.cacheManager.ClearPropertyCache(ctx, productID, deviceName)
	return err
}
func GetArrayID(id string, num int) string {
	return fmt.Sprintf("%s_%d", id, num)
}

// propertyDBIdentifiers 返回属性在普通时序库中的实际标识符。
func propertyDBIdentifiers(property schema.Property) []string {
	if property.Define.Type != schema.DataTypeArray {
		return []string{property.Identifier}
	}
	identifiers := make([]string, 0, cast.ToInt(property.Define.Max))
	for index := 0; index < cast.ToInt(property.Define.Max); index++ {
		identifiers = append(identifiers, GetArrayID(property.Identifier, index))
	}
	return identifiers
}

func (d *DeviceDataRepo) DeleteDeviceProperty(ctx context.Context, productID string, deviceName string, s []schema.Property) error {
	if len(s) > 0 {
		tableIdentifiers := make(map[string][]string)
		for _, v := range s {
			tableName := getTableName(v.Define)
			tableIdentifiers[tableName] = append(tableIdentifiers[tableName], propertyDBIdentifiers(v)...)
		}
		for tableName, identifiers := range tableIdentifiers {
			if len(identifiers) == 0 {
				continue
			}
			err := d.db.WithContext(ctx).Table(tableName).Where("product_id = ? and device_name = ? and identifier in ?", productID, deviceName, identifiers).Delete(&Property{}).Error
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
	err := d.cacheManager.ClearPropertyCache(ctx, productID, deviceName)
	return err
}
