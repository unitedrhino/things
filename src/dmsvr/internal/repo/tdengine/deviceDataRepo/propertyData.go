package deviceDataRepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"
	"time"
)

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *schema.Model, productID string, deviceName string, property *deviceData.PropertyData) error {
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

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string, params map[string]interface{}, timestamp time.Time) error {
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
	filter deviceData.FilterOpt) ([]*deviceData.PropertyData, error) {
	if err := filter.Check(); err != nil {
		return nil, err
	}

	var (
		err error
		sql sq.SelectBuilder
	)

	if filter.ArgFunc == "" {
		sql = sq.Select("*")
	} else {
		sql, err = d.GetPropertyArgFuncSelect(ctx, filter)
		if err != nil {
			return nil, err
		}
	}
	sql = sql.From(getPropertyStableName(filter.ProductID, filter.DataID)).
		Where("`device_name`=?", filter.DeviceName)
	sql = filter.Page.FmtSql(sql)
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
		retProperties = append(retProperties, ToPropertyData(filter.DataID, v))
	}
	return retProperties, err
}

func (d *DeviceDataRepo) GetPropertyArgFuncSelect(
	ctx context.Context,
	filter deviceData.FilterOpt) (sq.SelectBuilder, error) {
	schemaModel, err := d.getSchemaModel(ctx, filter.ProductID)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	p, ok := schemaModel.Property[filter.DataID]
	if !ok {
		return sq.SelectBuilder{}, errors.Parameter.AddMsgf("dataID:%s not find", filter.DataID)
	}
	var (
		sql sq.SelectBuilder
	)

	if p.Define.Type == schema.STRUCT {
		sql = sq.Select(getSpecsColumnWithArgFunc(p.Define.Specs, filter.ArgFunc))
	} else {
		sql = sq.Select(fmt.Sprintf("%s(`param`) as `param`", filter.ArgFunc))
	}
	if filter.Interval != 0 {
		sql = sql.Suffix("INTERVAL(?a)", filter.Interval)
	}
	return sql, nil
}

func (d *DeviceDataRepo) GetPropertyCountByID(
	ctx context.Context,
	filter deviceData.FilterOpt) (int64, error) {

	sqlData := sq.Select("count(1)").From(getPropertyStableName(filter.ProductID, filter.DataID)).
		Where("`device_name`=?", filter.DeviceName)
	sqlData = filter.Page.FmtWhere(sqlData)
	sqlStr, value, err := sqlData.ToSql()
	if err != nil {
		return 0, err
	}
	row := d.t.QueryRow(sqlStr, value...)
	var total int64
	err = row.Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
