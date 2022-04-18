package tdengine

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store/TDengine"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

const (
	PROPERTY_TYPE = "property_type"
)

type DeviceDataRepo struct {
	t *TDengine.Td
}

func NewDeviceDataRepo(dataSource string) *DeviceDataRepo {
	td, err := TDengine.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &DeviceDataRepo{t: td}
}

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

// GenParams 返回占位符?,?,?,? 参数id名:aa,bbb,ccc 参数值列表
func (d *DeviceDataRepo) GenParams(params map[string]interface{}) (string, string, []interface{}, error) {
	if len(params) == 0 {
		//使用这个函数前必须要判断参数的个数是否大于0
		panic("DeviceDataRepo|GenParams|params num == 0")
	}
	var (
		paramPlaceholder = strings.Repeat("?,", len(params))
		paramValList     []interface{} //参数值列表
		paramIds         []string
	)
	//将最后一个?去除
	paramPlaceholder = paramPlaceholder[:len(paramPlaceholder)-1]
	for k, v := range params {
		paramIds = append(paramIds, "`"+k+"`")
		if _, ok := v.([]interface{}); !ok {
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
