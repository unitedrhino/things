package schemaDataRepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *schema.Property, productID string, deviceName string,
	property *msgThing.Param, timestamp time.Time, optional msgThing.Optional) error {
	sql, args, err := d.GenInsertPropertySql(ctx, t, productID, deviceName, property, timestamp, optional)
	if err != nil {
		return err
	}
	d.t.AsyncInsert(sql, args...)
	return nil
}

func (d *DeviceDataRepo) GenInsertPropertySql(ctx context.Context, p *schema.Property, productID string, deviceName string,
	property *msgThing.Param, timestamp time.Time, optional msgThing.Optional) (sql string, args []any, err error) {
	var ars = map[string]any{}
	switch property.Define.Type {
	case schema.DataTypeArray:
		genArrSql := func(Identifier string, num int, v any) error {
			id := GetArrayID(Identifier, num)
			ars[schema.GenArray(Identifier, num)] = v
			switch vv := v.(type) {
			case map[string]any:
				paramPlaceholder, paramIds, paramValList, err := stores.GenParams(vv)
				if err != nil {
					return err
				}
				sql += fmt.Sprintf(" %s using %s tags('%s','%s',%d,'%s') (`ts`, %s) values (?,%s) ",
					d.GetPropertyTableName(productID, deviceName, id),
					d.GetPropertyStableName(p, productID, Identifier), productID, deviceName, num, p.Define.Type,
					paramIds, paramPlaceholder)
				args = append([]any{timestamp}, paramValList...)
			default:
				sql += fmt.Sprintf(" %s using %s tags('%s','%s',%d,'%s')(`ts`, `param`) values (?,?) ",
					d.GetPropertyTableName(productID, deviceName, id),
					d.GetPropertyStableName(p, productID, Identifier),
					productID, deviceName, num, p.Define.Type)
				args = append(args, timestamp, vv)
			}
			return nil
		}

		switch val := property.Value.(type) {
		case []any: //这种是数组的所有值一起上传的
			for i, v := range val {
				err := genArrSql(property.Identifier, i, v)
				if err != nil {
					return "", nil, err
				}
			}
		default:
			Identifier, num, ok := schema.GetArray(property.Identifier)
			if !ok {
				return "", nil, errors.Parameter.AddDetail("不是数组")
			}
			err := genArrSql(Identifier, num, val)
			if err != nil {
				return "", nil, err
			}
		}
	default:
		ars[property.Identifier] = property.Value
		switch property.Value.(type) {
		case map[string]any:
			paramPlaceholder, paramIds, paramValList, err := stores.GenParams(property.Value.(map[string]any))
			if err != nil {
				return "", nil, err
			}
			sql = fmt.Sprintf(" %s using %s tags('%s','%s','%s') (`ts`, %s) values (?,%s) ",
				d.GetPropertyTableName(productID, deviceName, property.Identifier),
				d.GetPropertyStableName(p, productID, property.Identifier), productID, deviceName, p.Define.Type,
				paramIds, paramPlaceholder)
			args = append([]any{timestamp}, paramValList...)
		default:
			var (
				param = property.Value
			)
			sql = fmt.Sprintf(" %s using %s tags('%s','%s','%s')(`ts`, `param`) values (?,?) ",
				d.GetPropertyTableName(productID, deviceName, property.Identifier),
				d.GetPropertyStableName(p, productID, property.Identifier),
				productID, deviceName, p.Define.Type)
			args = append(args, timestamp, param)
		}
	}
	f := func(ctx context.Context) {
		log := logx.WithContext(ctx)
		for k, v := range ars {
			var data = msgThing.PropertyData{
				Identifier: k,
				Param:      v,
				TimeStamp:  timestamp,
			}
			data.Fmt()
			err = d.kv.Hset(d.genRedisPropertyKey(productID, deviceName), k, data.String())
			if err != nil {
				log.Error(err)
			}
			retStr, err := d.kv.Hget(d.genRedisPropertyFirstKey(productID, deviceName), k)
			if err != nil && !errors.Cmp(stores.ErrFmt(err), errors.NotFind) {
				log.Error(err)
				continue
			}
			if retStr != "" {
				var ret msgThing.PropertyData
				err = json.Unmarshal([]byte(retStr), &ret)
				if err != nil {
					log.Error(err)
				} else if msgThing.IsParamValEq(&p.Define, v, ret.Param) { //相等不记录
					continue
				}
			}

			//到这里都是不相等或者之前没有记录的
			err = d.kv.Hset(d.genRedisPropertyFirstKey(productID, deviceName), k, data.String())
			if err != nil {
				log.Error(err)
			}
		}
	}
	if !optional.Sync {
		ctxs.GoNewCtx(ctx, f)
	} else {
		f(ctx)
	}

	return
}

func (d *DeviceDataRepo) genRedisPropertyFirstKey(productID string, deviceName string) string {
	return fmt.Sprintf("device:thing:property:first:%s:%s", productID, deviceName)
}

func (d *DeviceDataRepo) genRedisPropertyKey(productID string, deviceName string) string {
	return fmt.Sprintf("device:thing:property:last:%s:%s", productID, deviceName)
}
func (d *DeviceDataRepo) GetLatestPropertyDataByID(ctx context.Context, p *schema.Property, filter msgThing.LatestFilter) (*msgThing.PropertyData, error) {
	retStr, err := d.kv.HgetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName), filter.DataID)
	if err != nil && !errors.Cmp(stores.ErrFmt(err), errors.NotFind) {
		logx.WithContext(ctx).Error(err)
	}
	if retStr != "" {
		var ret msgThing.PropertyData
		err = json.Unmarshal([]byte(retStr), &ret)
		if err == nil {
			return &ret, nil
		}
	}
	//如果缓存里没有查到,需要从db里查
	dds, err := d.GetPropertyDataByID(ctx, p,
		msgThing.FilterOpt{
			Page:        def.PageInfo2{Size: 1},
			ProductID:   filter.ProductID,
			DeviceNames: []string{filter.DeviceName},
			DataID:      filter.DataID,
			Order:       stores.OrderDesc})
	if len(dds) == 0 || err != nil {
		return nil, err
	}
	d.kv.HsetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName), filter.DataID, dds[0].String())
	return dds[0], nil

}

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string,
	params map[string]msgThing.Param, timestamp time.Time, optional msgThing.Optional) error {
	var startTime = time.Now()
	defer func() {
		logx.WithContext(ctx).WithDuration(time.Now().Sub(startTime)).
			Infof("DeviceDataRepo.InsertPropertiesData")
	}()
	for identifier, param := range params {
		p := t.Property[param.Identifier]
		//入库
		param.Identifier = identifier
		sql1, args1, err := d.GenInsertPropertySql(ctx, p, productID, deviceName, &param, timestamp, optional)
		if err != nil {
			return errors.Database.AddDetailf(
				"DeviceDataRepo.InsertPropertiesData.InsertPropertyData identifier:%v param:%v err:%v",
				identifier, param, err)
		}
		if !optional.OnlyCache {
			d.t.AsyncInsert(sql1, args1...)
		}
	}
	return nil
}

func (d *DeviceDataRepo) GetPropertyDataByID(
	ctx context.Context, p *schema.Property,
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
		if filter.Order != stores.OrderAsc {
			sql = sql.OrderBy("`ts` desc")
		}
	} else {
		sql, err = d.getPropertyArgFuncSelect(ctx, p, filter)
		if err != nil {
			return nil, err
		}
		filter.Page.Size = 0
	}
	dataID := filter.DataID
	id, num, ok := schema.GetArray(filter.DataID)
	if ok {
		dataID = id
		sql = sql.Where("`_num`=?", num)
	}
	sql = sql.From(d.GetPropertyStableName(p, filter.ProductID, dataID))
	sql = d.fillFilter(sql, filter)
	sql = filter.Page.FmtSql(sql)

	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.QueryContext(ctx, sqlStr, value...)
	if err != nil {
		return nil, errors.Fmt(err).AddDetailf("sql:%v", sqlStr)
	}
	var datas []map[string]any
	stores.Scan(rows, &datas)
	retProperties := make([]*msgThing.PropertyData, 0, len(datas))
	for _, v := range datas {
		retProperties = append(retProperties, ToPropertyData(filter.DataID, p, v))
	}
	return retProperties, err
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect(
	ctx context.Context, p *schema.Property,
	filter msgThing.FilterOpt) (sq.SelectBuilder, error) {
	var (
		sql sq.SelectBuilder
	)

	if p.Define.Type == schema.DataTypeStruct {
		sql = sq.Select("FIRST(`ts`) AS ts,`device_name` ", d.GetSpecsColumnWithArgFunc(p.Define.Specs, filter.ArgFunc))
	} else {
		sql = sq.Select("FIRST(`ts`) AS ts,`device_name`  ", fmt.Sprintf("%s(`param`) as param", filter.ArgFunc))
	}
	if filter.Interval != 0 {
		var unit = filter.IntervalUnit
		if unit == "" {
			unit = "a"
		}
		sql = sql.Interval("?"+string(unit), filter.Interval)
	}
	if len(filter.Fill) > 0 {
		sql = sql.Fill(filter.Fill)
	}
	if filter.PartitionBy != "" {
		sql = sql.GroupBy(utils.CamelCaseToUdnderscore(filter.PartitionBy))
	}
	return sql, nil
}

func (d *DeviceDataRepo) fillFilter(
	sql sq.SelectBuilder, filter msgThing.FilterOpt) sq.SelectBuilder {
	if len(filter.DeviceNames) != 0 {
		sql = sql.Where(fmt.Sprintf("`device_name` in (%v)", stores.ArrayToSql(filter.DeviceNames)))
	}

	if len(filter.ProductIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`product_id` in (%v)", stores.ArrayToSql(filter.ProductIDs)))

	} else if filter.ProductID != "" {
		sql = sql.Where("`product_id` = ?", filter.ProductID)
	}

	//if filter.DeviceName != "" {
	//	sql = sql.Where("device_name=?", filter.DeviceName)
	//}
	return sql
}

func (d *DeviceDataRepo) GetPropertyCountByID(
	ctx context.Context, p *schema.Property,
	filter msgThing.FilterOpt) (int64, error) {
	sqlData := sq.Select("count(1)")
	dataID := filter.DataID
	id, num, ok := schema.GetArray(filter.DataID)
	if ok {
		dataID = id
		sqlData = sqlData.Where("`_num`=?", num)
	}
	sqlData = sqlData.From(d.GetPropertyStableName(p, filter.ProductID, dataID))
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
