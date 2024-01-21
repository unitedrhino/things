package schemaDataRepo

import (
	"fmt"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/stores"
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

func (S *SchemaStore) GetPropertyStableName(productID, identifier string) string {
	if productID != "" {
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
		tables = append(tables, S.GetPropertyStableName(productID, v.Identifier))
	}
	return
}
