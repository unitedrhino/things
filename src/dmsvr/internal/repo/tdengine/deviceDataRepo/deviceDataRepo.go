package deviceDataRepo

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	schema2 "github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

const (
	PROPERTY_TYPE = "property_type"
)

type DeviceDataRepo struct {
	t              *clients.Td
	getSchemaModel schema2.GetSchemaModel
}

func NewDeviceDataRepo(dataSource string, getSchemaModel schema2.GetSchemaModel) *DeviceDataRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td, getSchemaModel: getSchemaModel}
}

func getSpecsCreateColumn(s schema2.Specs) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("`%s` %s", v.ID, getTdType(v.DataType)))
	}
	return strings.Join(column, ",")
}

func getSpecsColumnWithArgFunc(s schema2.Specs, argFunc string) string {
	var column []string
	for _, v := range s {
		column = append(column, fmt.Sprintf("%s(`%s`) as %s", argFunc, v.ID, v.ID))
	}
	return strings.Join(column, ",")
}

func getTdType(define schema2.Define) string {
	switch define.Type {
	case schema2.BOOL:
		return "BOOL"
	case schema2.INT:
		return "BIGINT"
	case schema2.STRING:
		return "BINARY(5000)"
	case schema2.STRUCT:
		return "BINARY(5000)"
	case schema2.FLOAT:
		return "DOUBLE"
	case schema2.TIMESTAMP:
		return "TIMESTAMP"
	case schema2.ARRAY:
		return "BINARY(5000)"
	case schema2.ENUM:
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
func (d *DeviceDataRepo) GenParams(params map[string]any) (string, string, []any, error) {
	if len(params) == 0 {
		//使用这个函数前必须要判断参数的个数是否大于0
		panic("DeviceDataRepo|GenParams|params num == 0")
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
	t *schema2.Model,
	productID string,
	deviceName string) (tables []string) {
	for _, v := range t.Properties {
		tables = append(tables, getPropertyTableName(productID, deviceName, v.ID))
	}
	tables = append(tables, getEventTableName(productID, deviceName))
	return
}

func getStableNameList(
	t *schema2.Model,
	productID string) (tables []string) {
	for _, v := range t.Properties {
		tables = append(tables, getPropertyStableName(productID, v.ID))
	}
	tables = append(tables, getEventStableName(productID))
	return
}
