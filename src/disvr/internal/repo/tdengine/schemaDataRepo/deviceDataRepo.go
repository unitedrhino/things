package schemaDataRepo

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

const (
	PROPERTY_TYPE = "property_type"
)

type SchemaDataRepo struct {
	t              *clients.Td
	getSchemaModel schema.GetSchemaModel
}

func NewSchemaDataRepo(dataSource string, getSchemaModel schema.GetSchemaModel) *SchemaDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &SchemaDataRepo{t: td, getSchemaModel: getSchemaModel}
}

func getSpecsCreateColumn(s schema.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("`%s` %s", v.ID, getTdType(v.DataType)))
	}
	return strings.Join(column, ",")
}

func getSpecsColumnWithArgFunc(s schema.Specs, argFunc string) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("%s(`%s`) as %s", argFunc, v.ID, v.ID))
	}
	return strings.Join(column, ",")
}

func getTdType(define schema.Define) string {
	switch define.Type {
	case schema.BOOL:
		return "BOOL"
	case schema.INT:
		return "BIGINT"
	case schema.STRING:
		return "BINARY(5000)"
	case schema.STRUCT:
		return "BINARY(5000)"
	case schema.FLOAT:
		return "DOUBLE"
	case schema.TIMESTAMP:
		return "TIMESTAMP"
	case schema.ARRAY:
		return "BINARY(5000)"
	case schema.ENUM:
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

func getPropertyTableName(productID, deviceName, id string) string {
	return fmt.Sprintf("`device_property_%s_%s_%s`", productID, deviceName, id)
}
func getEventTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_event_%s_%s`", productID, deviceName)
}

// GenParams 返回占位符?,?,?,? 参数id名:aa,bbb,ccc 参数值列表
func (d *SchemaDataRepo) GenParams(params map[string]any) (string, string, []any, error) {
	if len(params) == 0 {
		//使用这个函数前必须要判断参数的个数是否大于0
		panic("SchemaDataRepo|GenParams|params num == 0")
	}
	var (
		paramPlaceholder = strings.Repeat("?,", len(params))
		paramValList     []any //参数值列表
		paramIds         []string
	)
	//将最后一个?去除
	paramPlaceholder = paramPlaceholder[:len(paramPlaceholder)-1]
	for k, v := range params {
		paramIds = append(paramIds, "`"+k+"`")
		if _, ok := v.([]any); !ok {
			paramValList = append(paramValList, v)
		} else { //如果是数组类型,需要序列化为json
			param, err := json.Marshal(v)
			if err != nil {
				return "", "", nil, errors.System.AddDetail("param json parse failure")
			}
			paramValList = append(paramValList, param)
		}
	}
	return paramPlaceholder, strings.Join(paramIds, ","), paramValList, nil
}

func getTableNameList(
	t *schema.Model,
	productID string,
	deviceName string) (tables []string) {
	for _, v := range t.Properties {
		tables = append(tables, getPropertyTableName(productID, deviceName, v.ID))
	}
	tables = append(tables, getEventTableName(productID, deviceName))
	return
}

func getStableNameList(
	t *schema.Model,
	productID string) (tables []string) {
	for _, v := range t.Properties {
		tables = append(tables, getPropertyStableName(productID, v.ID))
	}
	tables = append(tables, getEventStableName(productID))
	return
}
