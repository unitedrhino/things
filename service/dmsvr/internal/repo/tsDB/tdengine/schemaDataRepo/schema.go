package schemaDataRepo

import (
	"fmt"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"strings"
)

type SchemaStore struct {
}

func (S *SchemaStore) GetSpecsCreateColumn(s schema.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("`%s` %s", v.Identifier, tdengine.GetTdType(v.DataType)))
	}
	return strings.Join(column, ",")
}

func (S *SchemaStore) GetSpecsColumnWithArgFunc(s schema.Specs, argFunc string) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("%s(`%s`) as %s", argFunc, v.Identifier, v.Identifier))
	}
	return strings.Join(column, ",")
}

func (S *SchemaStore) GetPropertyStableName(p *schema.Property, productID, identifier string) string {
	if p != nil && p.Tag == schema.TagCustom && productID != "" {
		return fmt.Sprintf("`model_custom_property_%s_%s`", productID, identifier)
	}
	if p != nil && p.Tag == schema.TagDevice {
		switch p.Define.Type {
		case schema.DataTypeBool:
			return S.GetDeviceStableBoolName()
		case schema.DataTypeInt:
			return S.GetDeviceStableIntName()
		case schema.DataTypeString:
			return S.GetDeviceStableStringName()
		case schema.DataTypeFloat:
			return S.GetDeviceStableFloatName()
		case schema.DataTypeTimestamp:
			return S.GetDeviceStableTimestampName()
		case schema.DataTypeEnum:
			return S.GetDeviceStableEnumName()
		default:
			return ""
		}
	}
	return fmt.Sprintf("`model_common_property_%s`", identifier)
}

var defaultTags = "`product_id` ,`device_name`,`tenant_code`,`project_id`,`area_id`,`area_id_path` "
var defaultTagDef = "`product_id` BINARY(50),`device_name` BINARY(50),`" + PropertyType + "` BINARY(50)," +
	" `tenant_code`  BINARY(50),`project_id` BIGINT,`area_id` BIGINT,`area_id_path`  BINARY(50) "

func (S *SchemaStore) GetEventStableName() string {
	return fmt.Sprintf("`model_common_event`")
}

func (S *SchemaStore) GetPropertyTableName(productID, deviceName, identifier string) string {
	return fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, identifier)
}

var DeviceStables []string

func init() {
	s := SchemaStore{}
	DeviceStables = append(DeviceStables,
		s.GetDeviceStableBoolName(),
		s.GetDeviceStableIntName(),
		s.GetDeviceStableEnumName(),
		s.GetDeviceStableTimestampName(),
		s.GetDeviceStableFloatName(),
		s.GetDeviceStableStringName())
}

func (S *SchemaStore) GetDeviceStableBoolName() string {
	return fmt.Sprintf("`model_device_property_bool`")
}

func (S *SchemaStore) GetDeviceStableIntName() string {
	return fmt.Sprintf("`model_device_property_int`")
}

func (S *SchemaStore) GetDeviceStableEnumName() string {
	return fmt.Sprintf("`model_device_property_enum`")
}

func (S *SchemaStore) GetDeviceStableTimestampName() string {
	return fmt.Sprintf("`model_device_property_timestamp`")
}

func (S *SchemaStore) GetDeviceStableFloatName() string {
	return fmt.Sprintf("`model_device_property_float`")
}

func (S *SchemaStore) GetDeviceStableStringName() string {
	return fmt.Sprintf("`model_device_property_string`")
}

func (S *SchemaStore) GetPropertyTableNames(productID, deviceName string, p *schema.Property) (ret []string) {
	switch p.Define.Type {
	case schema.DataTypeArray:
		for i := 0; i < cast.ToInt(p.Define.Max); i++ {
			ret = append(ret, fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, GetArrayID(p.Identifier, i)))
		}
	default:
		return []string{fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, p.Identifier)}
	}
	return []string{}
}
func (S *SchemaStore) GetEventTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_event_%s_%s`", productID, deviceName)
}

func (S *SchemaStore) GetTableNameList(
	t *schema.Model,
	productID string,
	deviceName string) (tables []string) {
	for _, v := range t.Property {
		if v.Define.Type == schema.DataTypeArray {
			for i := 0; i < cast.ToInt(v.Define.Max); i++ {
				tables = append(tables, S.GetPropertyTableName(productID, deviceName, GetArrayID(v.Identifier, i)))
			}
			continue
		}
		tables = append(tables, S.GetPropertyTableName(productID, deviceName, v.Identifier))
	}
	tables = append(tables, S.GetEventTableName(productID, deviceName))
	return
}

func (S *SchemaStore) GetStableNameList(
	t *schema.Model,
	productID string) (tables []string) {
	if t.Property == nil {
		return []string{}
	}
	for _, v := range t.Property {
		if v.Tag == schema.TagCustom {
			tables = append(tables, S.GetPropertyStableName(v, productID, v.Identifier))
		}
	}
	return
}
