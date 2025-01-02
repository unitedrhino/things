package schemaDataRepo

import (
	"fmt"
	"gitee.com/unitedrhino/share/domain/schema"
)

var TableNames []string

var types = []schema.Define{
	{Type: schema.DataTypeBool},
	{Type: schema.DataTypeInt},
	{Type: schema.DataTypeString},
	{Type: schema.DataTypeTimestamp},
	{Type: schema.DataTypeEnum},
	{Type: schema.DataTypeFloat},
	{Type: schema.DataTypeArray, ArrayInfo: &schema.Define{Type: schema.DataTypeBool}},
	{Type: schema.DataTypeArray, ArrayInfo: &schema.Define{Type: schema.DataTypeInt}},
	{Type: schema.DataTypeArray, ArrayInfo: &schema.Define{Type: schema.DataTypeString}},
	{Type: schema.DataTypeArray, ArrayInfo: &schema.Define{Type: schema.DataTypeTimestamp}},
	{Type: schema.DataTypeArray, ArrayInfo: &schema.Define{Type: schema.DataTypeEnum}},
	{Type: schema.DataTypeArray, ArrayInfo: &schema.Define{Type: schema.DataTypeFloat}},
}

func init() {
	for _, t := range types {
		TableNames = append(TableNames, getTableName(t))
	}
}

func getModel(p schema.Define) interface{} {
	var isArray bool
	if p.Type == schema.DataTypeArray {
		isArray = true
		p = *p.ArrayInfo
	}
	switch p.Type {
	case schema.DataTypeBool:
		if isArray {
			return PropertyBoolArray{}
		}
		return PropertyBool{}
	case schema.DataTypeInt:
		if isArray {
			return PropertyIntArray{}
		}
		return PropertyInt{}
	case schema.DataTypeString:
		if isArray {
			return PropertyStringArray{}
		}
		return PropertyString{}
	case schema.DataTypeStruct:
		if isArray {
			return PropertyStructArray{}
		}
		return PropertyStruct{}
	case schema.DataTypeFloat:
		if isArray {
			return PropertyFloatArray{}
		}
		return PropertyFloat{}
	case schema.DataTypeTimestamp:
		if isArray {
			return PropertyTimestampArray{}
		}
		return PropertyTimestamp{}
	case schema.DataTypeEnum:
		if isArray {
			return PropertyEnumArray{}
		}
		return PropertyEnum{}
	}
	return nil
}

func getTableName(p schema.Define) string {
	var isArray bool
	if p.Type == schema.DataTypeArray {
		isArray = true
		p = *p.ArrayInfo
	}
	tableName := fmt.Sprintf("dm_model_property_%s", p.Type)
	if isArray {
		tableName += "_array"
	}
	return tableName
}
