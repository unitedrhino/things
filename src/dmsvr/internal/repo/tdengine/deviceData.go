package tdengine

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
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
	//todo 后续重构为一条sql插入 向多个表插入记录 参考:https://www.taosdata.com/docs/cn/v2.0/taos-sql#management
	for id, param := range params {
		err := d.InsertPropertyData(ctx, productID, deviceName, &deviceTemplate.PropertyData{
			ID:        id,
			Param:     param,
			TimeStamp: timestamp,
		})
		if err != nil {
			return errors.Database.AddDetailf("DeviceDataRepo|InsertPropertiesData|InsertPropertyData id:%v param:%v err:%v",
				id, param, err)
		}
	}
	return nil
}

func (d *DeviceDataRepo) GetEventDataWithID(ctx context.Context, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*deviceTemplate.EventData, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceDataRepo) GetPropertyDataByID(ctx context.Context, productID string, deviceName string, dataID string, page def.PageInfo2) ([]*deviceTemplate.PropertyData, error) {
	//select * from model_property_23fipsijpsk_wifi_info where device_name='test5' and ts>'2022-04-14 22:22:30'  limit 10
	rows, err := d.t.Query("select * from model_property_23fipsijpsk_wifi_info where device_name='test5' and ts>'2022-04-14 22:22:30'  limit 10")
	if err != nil {
		return nil, err
	}
	var datas []map[string]interface{}
	store.Scan(rows, &datas)
	fmt.Println(datas)
	//for rows.Next() {
	//	var params = make([]interface{}, 3)
	//	for index, _ := range params { //为每一列初始化一个指针
	//		var a interface{}
	//		params[index] = &a
	//	}
	//	err := rows.Scan(params...)
	//	fmt.Println(err)
	//	columnType, err := rows.ColumnTypes()
	//	fmt.Println(columnType, err)
	//	columns, err := rows.Columns()
	//	fmt.Println(columns, err)
	//}
	return nil, err
}

func GetRows(rows *sql.Rows) []map[string]interface{} {
	defer rows.Close()
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	return list
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
