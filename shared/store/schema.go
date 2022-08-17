package store

import (
	"fmt"
	"github.com/i-Things/things/shared/domain/schema"
	"strings"
)

type SchemaStore struct {
}

func (S *SchemaStore) GetSpecsCreateColumn(s schema.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("`%s` %s", v.ID, GetTdType(v.DataType)))
	}
	return strings.Join(column, ",")
}

func (S *SchemaStore) GetSpecsColumnWithArgFunc(s schema.Specs, argFunc string) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("%s(`%s`) as %s", argFunc, v.ID, v.ID))
	}
	return strings.Join(column, ",")
}

func (S *SchemaStore) GetPropertyStableName(productID, id string) string {
	return fmt.Sprintf("`model_property_%s_%s`", productID, id)
}
func (S *SchemaStore) GetEventStableName(productID string) string {
	return fmt.Sprintf("`model_event_%s`", productID)
}

func (S *SchemaStore) GetPropertyTableName(productID, deviceName, id string) string {
	return fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, id)
}
func (S *SchemaStore) GetEventTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_event_%s_%s`", productID, deviceName)
}

func (S *SchemaStore) GetTableNameList(
	t *schema.Model,
	productID string,
	deviceName string) (tables []string) {
	for _, v := range t.Properties {
		tables = append(tables, S.GetPropertyTableName(productID, deviceName, v.ID))
	}
	tables = append(tables, S.GetEventTableName(productID, deviceName))
	return
}

func (S *SchemaStore) GetStableNameList(
	t *schema.Model,
	productID string) (tables []string) {
	for _, v := range t.Properties {
		tables = append(tables, S.GetPropertyStableName(productID, v.ID))
	}
	tables = append(tables, S.GetEventStableName(productID))
	return
}
