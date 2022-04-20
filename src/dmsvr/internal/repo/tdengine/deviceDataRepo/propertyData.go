package deviceDataRepo

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"
	"github.com/i-Things/things/src/dmsvr/internal/domain/templateModel"
	"time"
)

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *templateModel.Template, productID string, deviceName string, property *deviceData.PropertyData) error {
	switch property.Param.(type) {
	case map[string]interface{}: //结构体类型
		paramPlaceholder, paramIds, paramValList, err := d.GenParams(property.Param.(map[string]interface{}))
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("insert into %s using %s tags('%s','%s') (ts, %s) values (?,%s);",
			getPropertyTableName(productID, deviceName, property.ID),
			getPropertyStableName(productID, property.ID), deviceName, t.Property[property.ID].Define.Type,
			paramIds, paramPlaceholder)
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

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, t *templateModel.Template, productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error {
	//todo 后续重构为一条sql插入 向多个表插入记录 参考:https://www.taosdata.com/docs/cn/v2.0/taos-sql#management
	for id, param := range params {
		err := d.InsertPropertyData(ctx, t, productID, deviceName, &deviceData.PropertyData{
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

func (d *DeviceDataRepo) GetPropertyDataByID(
	ctx context.Context,
	productID string,
	deviceName string,
	dataID string,
	page def.PageInfo2) ([]*deviceData.PropertyData, error) {

	sql := sq.Select("*").From(getPropertyStableName(productID, dataID)).
		Where("`device_name`=?", deviceName).OrderBy("`ts` desc")
	sql = page.FmtSql(sql)
	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.Query(sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]interface{}
	store.Scan(rows, &datas)
	retProperties := make([]*deviceData.PropertyData, 0, len(datas))
	for _, v := range datas {
		retProperties = append(retProperties, ToPropertyData(dataID, v))
	}
	return retProperties, err
}
