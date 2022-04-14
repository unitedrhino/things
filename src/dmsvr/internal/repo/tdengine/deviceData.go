package tdengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store/TDengine"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
	"time"
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

func (d *DeviceDataRepo) InsertEventData(ctx context.Context, productID string,
	deviceName string, event *deviceTemplate.EventData) error {
	param, err := json.Marshal(event.Params)
	if err != nil {
		return errors.System.AddDetail("param json parse failure")
	}
	sql := fmt.Sprintf("insert into %s (ts,event_id,event_type, param) values (?,?,?,?);", getEventTableName(productID, deviceName))
	if _, err := d.t.Exec(sql, event.TimeStamp, event.ID, event.Type, param); err != nil {
		return err
	}
	return nil
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
		paramIds = append(paramIds, k)
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

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, productID string, deviceName string, property *deviceTemplate.PropertyData) error {
	switch property.Param.(type) {
	case map[string]interface{}: //结构体类型
		paramPlaceholder, paramIds, paramValList, err := d.GenParams(property.Param.(map[string]interface{}))
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("insert into %s (ts, %s) values (?,%s);",
			getPropertyTableName(productID, deviceName, property.ID), paramIds, paramPlaceholder)
		param := append([]interface{}{property.TimeStamp}, paramValList...)
		if _, err := d.t.Exec(sql, param...); err != nil {
			return err
		}
	default:
		var (
			param = property.Param
			err   error
		)
		if _, ok := property.Param.([]interface{}); ok { //如果是数组类型,需要先序列化为json
			param, err = json.Marshal(property.Param)
			if err != nil {
				return errors.System.AddDetail("param json parse failure")
			}
		}
		sql := fmt.Sprintf("insert into %s (ts, param) values (?,?);", getPropertyTableName(productID, deviceName, property.ID))
		if _, err := d.t.Exec(sql, property.TimeStamp, param); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) GetEventDataWithID(ctx context.Context, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*deviceTemplate.EventData, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) GetPropertyDataByID(ctx context.Context, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*deviceTemplate.PropertyData, error) {
	//TODO implement me
	panic("implement me")
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
	return fmt.Sprintf("model_property_%s_%s", productID, id)
}
func getEventStableName(productID string) string {
	return fmt.Sprintf("model_event_%s", productID)
}

func getActionStableName(productID string) string {
	return fmt.Sprintf("model_action_%s", productID)
}

func getPropertyTableName(productID, deviceName, id string) string {
	return fmt.Sprintf("device_property_%s_%s_%s", productID, deviceName, id)
}
func getEventTableName(productID, deviceName string) string {
	return fmt.Sprintf("device_event_%s_%s", productID, deviceName)
}

func getActionTableName(productID, deviceName string) string {
	return fmt.Sprintf("device_action_%s_%s", productID, deviceName)
}
