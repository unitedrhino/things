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
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *schema.Model, productID string, deviceName string, property *msgThing.PropertyData) error {
	sql, args, err := d.GenInsertPropertySql(ctx, t, productID, deviceName, property)
	if err != nil {
		return err
	}
	d.t.AsyncInsert(sql, args)
	return nil
}

func (d *DeviceDataRepo) GenInsertPropertySql(ctx context.Context, t *schema.Model, productID string, deviceName string, property *msgThing.PropertyData) (sql string, args []any, err error) {
	switch property.Param.(type) {
	case map[string]any:
		paramPlaceholder, paramIds, paramValList, err := stores.GenParams(property.Param.(map[string]any))
		if err != nil {
			return "", nil, err
		}
		sql = fmt.Sprintf(" %s using %s tags('%s','%s') (`ts`, %s) values (?,%s) ",
			d.GetPropertyTableName(productID, deviceName, property.Identifier),
			d.GetPropertyStableName(productID, property.Identifier), deviceName, t.Property[property.Identifier].Define.Type,
			paramIds, paramPlaceholder)
		args = append([]any{property.TimeStamp}, paramValList...)
	default:
		var (
			param = property.Param
			err   error
		)
		if _, ok := property.Param.([]any); ok { //如果是数组类型,需要先序列化为json
			param, err = json.Marshal(property.Param)
			if err != nil {
				return "", nil, errors.System.AddDetail("param json parse failure")
			}
		}
		sql = fmt.Sprintf(" %s using %s tags('%s','%s')(`ts`, `param`) values (?,?) ",
			d.GetPropertyTableName(productID, deviceName, property.Identifier),
			d.GetPropertyStableName(productID, property.Identifier),
			deviceName, t.Property[property.Identifier].Define.Type)
		args = append(args, property.TimeStamp, param)
	}
	return
}

func (d *DeviceDataRepo) genRedisPropertyKey(productID string, deviceName, identifier string) string {
	return fmt.Sprintf("device:thing:property:%s:%s:%s", productID, deviceName, identifier)
}
func (d *DeviceDataRepo) GetLatestPropertyDataByID(ctx context.Context, filter msgThing.LatestFilter) (*msgThing.PropertyData, error) {
	retStr, err := d.kv.GetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName, filter.DataID))
	if err != nil {
		return nil, errors.Database.AddDetailf(
			"DeviceDataRepo.GetLatestPropertyDataByID.GetCtx filter:%v  err:%v",
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

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string, params map[string]any, timestamp time.Time) error {
	var startTime = time.Now()
	defer func() {
		logx.WithContext(ctx).WithDuration(time.Now().Sub(startTime)).
			Infof("DeviceDataRepo.InsertPropertiesData")
	}()
	for identifier, param := range params {
		data := msgThing.PropertyData{
			Identifier: identifier,
			Param:      param,
			TimeStamp:  timestamp,
		}
		//缓存
		err := d.kv.SetCtx(ctx, d.genRedisPropertyKey(productID, deviceName, identifier), data.String())
		if err != nil {
			return errors.Database.AddDetailf(
				"DeviceDataRepo.InsertPropertiesData.SetCtx identifier:%v param:%v err:%v",
				identifier, param, err)
		}
		//入库
		sql1, args1, err := d.GenInsertPropertySql(ctx, t, productID, deviceName, &data)
		if err != nil {
			return errors.Database.AddDetailf(
				"DeviceDataRepo.InsertPropertiesData.InsertPropertyData identifier:%v param:%v err:%v",
				identifier, param, err)
		}
		d.t.AsyncInsert(sql1, args1)
	}
	return nil
}

func (d *DeviceDataRepo) GetPropertyDataByID(
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
			sql = sql.OrderBy("`ts` desc")
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
	stores.Scan(rows, &datas)
	retProperties := make([]*msgThing.PropertyData, 0, len(datas))
	for _, v := range datas {
		retProperties = append(retProperties, ToPropertyData(filter.DataID, v))
	}
	return retProperties, err
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect(
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
		sql = sq.Select("FIRST(`ts`) AS ts", d.GetSpecsColumnWithArgFunc(p.Define.Specs, filter.ArgFunc))
	} else {
		sql = sq.Select("FIRST(`ts`) AS ts", fmt.Sprintf("%s(`param`) as param", filter.ArgFunc))
	}
	if filter.Interval != 0 {
		sql = sql.Interval("?a", filter.Interval)
	}
	if len(filter.Fill) > 0 {
		sql = sql.Fill(filter.Fill)
	}
	return sql, nil
}

func (d *DeviceDataRepo) fillFilter(
	sql sq.SelectBuilder, filter msgThing.FilterOpt) sq.SelectBuilder {
	if len(filter.DeviceNames) != 0 {
		sql = sql.Where(fmt.Sprintf("device_name= (%v)", stores.ArrayToSql(filter.DeviceNames)))
	}
	return sql
}

func (d *DeviceDataRepo) GetPropertyCountByID(
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
