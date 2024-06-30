package schemaDataRepo

import (
	"fmt"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/stores"
	"github.com/spf13/cast"
	"strings"
)

type SchemaStore struct {
}

func (S *SchemaStore) GetSpecsCreateColumn(s schema.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("`%s` %s", v.Identifier, stores.GetTdType(v.DataType)))
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

func (S *SchemaStore) GetPropertyStableName(tag schema.Tag, productID, identifier string) string {
	if tag == schema.TagCustom && productID != "" {
		return fmt.Sprintf("`model_custom_property_%s_%s`", productID, identifier)
	}
	return fmt.Sprintf("`model_common_property_%s`", identifier)
}
func (S *SchemaStore) GetEventStableName() string {
	return fmt.Sprintf("`model_common_event`")
}

func (S *SchemaStore) GetPropertyTableName(productID, deviceName, identifier string) string {
	return fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, identifier)
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
			tables = append(tables, S.GetPropertyStableName(v.Tag, productID, v.Identifier))
		}
	}
	return
}
