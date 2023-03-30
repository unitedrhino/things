package schemaDataRepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"time"
)

func (d *SchemaDataRepo) InsertPropertyData(ctx context.Context, t *schema.Model, productID string, deviceName string, property *msgThing.PropertyData) error {
	switch property.Param.(type) {
	case map[string]any:
		paramPlaceholder, paramIds, paramValList, err := store.GenParams(property.Param.(map[string]any))
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("insert into %s using %s tags('%s','%s') (ts, %s) values (?,%s);",
			d.GetPropertyTableName(productID, deviceName, property.Identifier),
			d.GetPropertyStableName(productID, property.Identifier), deviceName, t.Property[property.Identifier].Define.Type,
			paramIds, paramPlaceholder)
		param := append([]any{property.TimeStamp}, paramValList...)
		if _, err := d.t.ExecContext(ctx, sql, param...); err != nil {
			return err
		}
	default:
		var (
			param = property.Param
			err   error
		)
		if _, ok := property.Param.([]any); ok { //如果是数组类型,需要先序列化为json
			param, err = json.Marshal(property.Param)
			if err != nil {
				return errors.System.AddDetail("param json parse failure")
			}
		}
		sql := fmt.Sprintf("insert into %s using %s tags('%s','%s')(ts, param) values (?,?);",
			d.GetPropertyTableName(productID, deviceName, property.Identifier),
			d.GetPropertyStableName(productID, property.Identifier),
			deviceName, t.Property[property.Identifier].Define.Type)
		if _, err := d.t.ExecContext(ctx, sql, property.TimeStamp, param); err != nil {
			return err
		}
	}
	return nil
}

func (d *SchemaDataRepo) genRedisPropertyKey(productID string, deviceName, identifier string) string {
	return fmt.Sprintf("device:thing:property:%s:%s:%s", productID, deviceName, identifier)
}
func (d *SchemaDataRepo) GetLatestPropertyDataByID(ctx context.Context, filter msgThing.LatestFilter) (*msgThing.PropertyData, error) {
	retStr, err := d.kv.GetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName, filter.DataID))
	if err != nil {
		return nil, errors.Database.AddDetailf(
			"SchemaDataRepo.GetLatestPropertyDataByID.GetCtx filter:%v  err:%v",
			filter, err)
	}
	if retStr == "" { //如果缓存里没有查到,需要从db里查
		dds, err := d.GetPropertyDataByID(ctx,
			msgThing.FilterOpt{
				Page:        def.PageInfo2{Size: 1},
				ProductID:   filter.ProductID,
				DeviceNames: []string{filter.DeviceName},
				DataID:      filter.DataID,
				Order:       def.OrderDesc})
		if len(dds) == 0 || err != nil {
			return nil, err
		}
		d.kv.SetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName, filter.DataID), dds[0].String())
		return dds[0], nil
	}
	var ret msgThing.PropertyData
	err = json.Unmarshal([]byte(retStr), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (d *SchemaDataRepo) InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string, params map[string]any, timestamp time.Time) error {
	//todo 后续重构为一条sql插入 向多个表插入记录 参考:https://www.taosdata.com/docs/cn/v2.0/taos-sql#management
	for identifier, param := range params {
		data := msgThing.PropertyData{
			Identifier: identifier,
			Param:      param,
			TimeStamp:  timestamp,
		}
		err := d.kv.SetCtx(ctx, d.genRedisPropertyKey(productID, deviceName, identifier), data.String())
		if err != nil {
			return errors.Database.AddDetailf(
				"SchemaDataRepo.InsertPropertiesData.SetCtx identifier:%v param:%v err:%v",
				identifier, param, err)
		}
		err = d.InsertPropertyData(ctx, t, productID, deviceName, &data)
		if err != nil {
			return errors.Database.AddDetailf(
				"SchemaDataRepo.InsertPropertiesData.InsertPropertyData identifier:%v param:%v err:%v",
				identifier, param, err)
		}
	}
	return nil
}

func (d *SchemaDataRepo) GetPropertyDataByID(
	ctx context.Context,
	filter msgThing.FilterOpt) ([]*msgThing.PropertyData, error) {
	if err := filter.Check(); err != nil {
		return nil, err
	}

	var (
		err error
		sql sq.SelectBuilder
	)

	if filter.ArgFunc == "" {
		sql = sq.Select("*")
		if filter.Order != def.OrderAes {
			sql = sql.OrderBy("ts desc")
		}
	} else {
		sql, err = d.getPropertyArgFuncSelect(ctx, filter)
		if err != nil {
			return nil, err
		}
		filter.Page.Size = 0
	}
	sql = sql.From(d.GetPropertyStableName(filter.ProductID, filter.DataID))
	sql = d.fillFilter(sql, filter)
	sql = filter.Page.FmtSql(sql)

	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.QueryContext(ctx, sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]any
	store.Scan(rows, &datas)
	retProperties := make([]*msgThing.PropertyData, 0, len(datas))
	for _, v := range datas {
		retProperties = append(retProperties, ToPropertyData(filter.DataID, v))
	}
	return retProperties, err
}

func (d *SchemaDataRepo) getPropertyArgFuncSelect(
	ctx context.Context,
	filter msgThing.FilterOpt) (sq.SelectBuilder, error) {
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

	if p.Define.Type == schema.DataTypeStruct {
		sql = sq.Select("FIRST(ts) AS `ts`", d.GetSpecsColumnWithArgFunc(p.Define.Specs, filter.ArgFunc))
	} else {
		sql = sq.Select("FIRST(ts) AS `ts`", fmt.Sprintf("%s(`param`) as `param`", filter.ArgFunc))
	}
	if filter.Interval != 0 {
		sql = sql.Interval("?a", filter.Interval)
	}
	if len(filter.Fill) > 0 {
		sql = sql.Fill(filter.Fill)
	}
	return sql, nil
}

func (d *SchemaDataRepo) fillFilter(
	sql sq.SelectBuilder, filter msgThing.FilterOpt) sq.SelectBuilder {
	if len(filter.DeviceNames) != 0 {
		sql = sql.Where(fmt.Sprintf("`deviceName` in (%v)", store.ArrayToSql(filter.DeviceNames)))
	}
	return sql
}

func (d *SchemaDataRepo) GetPropertyCountByID(
	ctx context.Context,
	filter msgThing.FilterOpt) (int64, error) {

	sqlData := sq.Select("count(1)").From(d.GetPropertyStableName(filter.ProductID, filter.DataID))
	sqlData = d.fillFilter(sqlData, filter)
	sqlData = filter.Page.FmtWhere(sqlData)
	sqlStr, value, err := sqlData.ToSql()
	if err != nil {
		return 0, err
	}
	row := d.t.QueryRowContext(ctx, sqlStr, value...)
	var total int64
	err = row.Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
