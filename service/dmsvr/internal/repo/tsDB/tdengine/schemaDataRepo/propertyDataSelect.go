package schemaDataRepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func (d *DeviceDataRepo) GetLatestPropertyDataByID(ctx context.Context, p *schema.Property, filter msgThing.LatestFilter) (*msgThing.PropertyData, error) {
	retStr, err := d.kv.HgetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName), filter.DataID)
	if err != nil && !errors.Cmp(stores.ErrFmt(err), errors.NotFind) {
		logx.WithContext(ctx).Error(err)
	}
	if retStr != "" {
		var ret msgThing.PropertyData
		err = json.Unmarshal([]byte(retStr), &ret)
		if err == nil {
			vv, er := msgThing.GetVal(&p.Define, ret.Param)
			if er == nil {
				ret.Param = vv
			}
			return &ret, nil
		}
	}
	//如果缓存里没有查到,需要从db里查
	dds, err := d.GetPropertyDataByID(ctx, p,
		msgThing.FilterOpt{
			Filter: msgThing.Filter{
				ProductID:   filter.ProductID,
				DeviceNames: []string{filter.DeviceName},
			},
			Page:   def.PageInfo2{Size: 1},
			DataID: filter.DataID,
			Order:  stores.OrderDesc})
	if len(dds) == 0 || err != nil {
		return nil, err
	}
	vv, er := msgThing.GetVal(&p.Define, dds[0].Param)
	if er == nil {
		dds[0].Param = vv
	}
	d.kv.HsetCtx(ctx, d.genRedisPropertyKey(filter.ProductID, filter.DeviceName), filter.DataID, dds[0].String())
	return dds[0], nil

}

//func (d *DeviceDataRepo) GetPropertyAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterOpt) ([]*msgThing.PropertyDatas, error) {
//	//TODO implement me
//	panic("implement me")
//}

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
	sql = d.fillFilter(sql, filter.Filter)
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
		retProperties = append(retProperties, d.ToPropertyData(ctx, filter.DataID, p, v))
	}
	return retProperties, err
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect(
	ctx context.Context, p *schema.Property,
	filter msgThing.FilterOpt) (sq.SelectBuilder, error) {
	var (
		sql sq.SelectBuilder
	)
	deviceName := ",`device_name` "
	partitionBy := utils.CamelCaseToUdnderscore(filter.PartitionBy)
	if !strings.Contains(partitionBy, "device_name") { //如果没有传partition by 会报错
		deviceName = ""
	}
	pb := partitionBy
	if partitionBy != "" {
		pb = "," + pb
	}
	ts := "FIRST(`ts`)  AS ts "
	if filter.Interval != 0 {
		ts = "_wstart AS ts "
	} else if filter.NoFirstTs {
		ts = "`ts` "
	}
	if p.Define.Type == schema.DataTypeStruct {
		sql = sq.Select(ts+deviceName+pb, d.GetSpecsColumnWithArgFunc(p.Define.Specs, filter.ArgFunc))
	} else {
		sql = sq.Select(ts+deviceName+pb, fmt.Sprintf("%s(`param`) as param", filter.ArgFunc))
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
		sql = sql.PartitionBys(partitionBy)
	}
	return sql, nil
}

func (d *DeviceDataRepo) fillFilter(
	sql sq.SelectBuilder, filter msgThing.Filter) sq.SelectBuilder {
	if len(filter.DeviceNames) != 0 {
		sql = sql.Where(fmt.Sprintf("`device_name` in (%v)", stores.ArrayToSql(filter.DeviceNames)))
	}

	if len(filter.ProductIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`product_id` in (%v)", stores.ArrayToSql(filter.ProductIDs)))

	} else if filter.ProductID != "" {
		sql = sql.Where("`product_id` = ?", filter.ProductID)
	}

	if filter.TenantCode != "" {
		sql = sql.Where("`tenant_code`=?", filter.TenantCode)
	}
	sql = tdengine.GroupFilter(sql, d.groupConfigs, filter.BelongGroup)

	if filter.ProjectID != 0 {
		sql = sql.Where("`project_id`=?", filter.ProjectID)
	}
	if filter.AreaID != 0 {
		sql = sql.Where("`area_id`=?", filter.AreaID)
	}
	if filter.AreaIDPath != "" {
		sql = sql.Where("`area_id_path` like ?", filter.AreaIDPath+"%")
	}
	if len(filter.AreaIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`area_id` in (%v)", stores.ArrayToSql(filter.AreaIDs)))
	}
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
	sqlData = d.fillFilter(sqlData, filter.Filter)
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
