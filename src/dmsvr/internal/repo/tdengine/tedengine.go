package tdengine

import (
	"fmt"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"strings"
)

const (
	PROPERTY_TYPE = "property_type"
)

func getSpecsColumn(s deviceTemplate.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("`%s` %s", v.ID, getTdType(v.DataType)))
	}
	return strings.Join(column, ",")
}

func getTdType(define deviceTemplate.Define) string {
	switch define.Type {
	case deviceTemplate.BOOL:
		return "BOOL"
	case deviceTemplate.INT:
		return "BIGINT"
	case deviceTemplate.STRING:
		return fmt.Sprintf("BINARY(%s)", define.Max)
	case deviceTemplate.STRUCT:
		return "BINARY(5000)"
	case deviceTemplate.FLOAT:
		return "DOUBLE"
	case deviceTemplate.TIMESTAMP:
		return "TIMESTAMP"
	case deviceTemplate.ARRAY:
		return "BINARY(5000)"
	case deviceTemplate.ENUM:
		return "SMALLINT"
	default:
		panic(fmt.Sprintf("%v not support", define.Type))
	}
}

func getPropertyStableName(productID, id string) string {
	return fmt.Sprintf("`model_property_%s_%s`", productID, id)
}
func getEventStableName(productID string) string {
	return fmt.Sprintf("`model_event_%s`", productID)
}

func getActionStableName(productID string) string {
	return fmt.Sprintf("`model_action_%s`", productID)
}

func getPropertyTableName(productID, deviceName, id string) string {
	return fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, id)
}
func getEventTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_event_%s_%s`", productID, deviceName)
}

func getActionTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_action_%s_%s`", productID, deviceName)
}
