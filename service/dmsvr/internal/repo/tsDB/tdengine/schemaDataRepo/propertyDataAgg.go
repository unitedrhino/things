package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type PropertyAggStu struct {
	NoFirstTs bool `json:"noFirstTs,optional"` //时间戳填充不填充最早的值,聚合模式使用
	Aggs      []PropertyAgg2
}
type PropertyAgg2 struct {
	P *schema.Property
	msgThing.PropertyAgg
}

func (d *DeviceDataRepo) GetPropertyAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterAggOpt) ([]*msgThing.PropertyData2, error) {
	var (
		err error
	)
	if len(filter.Aggs) == 0 {
		return nil, errors.Parameter.AddMsg("至少填写一个聚合函数")
	}
	for _, agg := range filter.Aggs { //todo 暂时不考虑数组类型
		_, ok := m.Property[agg.DataID]
		if !ok {
			return nil, errors.Parameter.AddMsgf("属性未定义:%v", agg.DataID)
		}
	}
	var retMap = map[string]msgThing.PropertyData2{}
	var page = def.PageInfo2{TimeStart: filter.TimeStart, TimeEnd: filter.TimeEnd}
	for _, agg := range filter.Aggs { //暂时不考虑数组类型
		p, ok := m.Property[agg.DataID]
		if !ok {
			return nil, errors.Parameter.AddMsgf("属性未定义:%v", agg.DataID)
		}
		sql, err := d.getPropertyArgFuncSelect2(ctx, p, agg, filter)
		if err != nil {
			logx.WithContext(ctx).Errorf(err.Error())
			return nil, err
		}
		dataID := agg.DataID
		id, num, ok := schema.GetArray(agg.DataID)
		if ok {
			dataID = id
			sql = sql.Where("`_num`=?", num)
		}
		sql = sql.From(d.GetPropertyStableName(p, filter.ProductID, dataID))
		sql = d.fillFilter(sql, filter.Filter)
		sql = page.FmtSql(sql)
		sqlStr, value, err := sql.ToSql()
		if err != nil {
			logx.WithContext(ctx).Errorf(err.Error())
			return nil, err
		}
		rows, err := d.t.QueryContext(ctx, sqlStr, value...)
		if err != nil {
			logx.WithContext(ctx).Errorf("sql:%v err:%v", sqlStr, err.Error())
			return nil, err
		}
		var datas []map[string]any
		stores.Scan(rows, &datas)
		retMap = d.ToPropertyData2(ctx, agg, m, datas, retMap)
	}

	return utils.MapVToSlice2(retMap), err
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect2(
	ctx context.Context, p *schema.Property, agg msgThing.PropertyAgg,
	filter msgThing.FilterAggOpt) (sq.SelectBuilder, error) {
	var (
		sql sq.SelectBuilder
	)
	partitionBy := utils.CamelCaseToUdnderscore(filter.PartitionBy)
	var selects = []string{"_wstart AS ts_window "}
	if partitionBy != "" {
		selects = append(selects, partitionBy)
	}
	if p.Define.Type == schema.DataTypeStruct { //todo 暂未支持
	} else {
		for _, argFunc := range agg.ArgFuncs {
			//pg的 timescale走视图优化
			if agg.NoFirstTs && utils.SliceIn(argFunc, "first", "last", "min", "max") {
				selects = append(selects, fmt.Sprintf(` %s(param) as %s_param,cols(%s(param),ts) as %s_ts `, argFunc, argFunc, argFunc, argFunc))
			} else {
				selects = append(selects, fmt.Sprintf(` %s(param) as %s_param `, argFunc, argFunc))
			}
		}
	}
	sql = sq.Select(selects...)
	if filter.Interval != 0 {
		var unit = filter.IntervalUnit
		if unit == "" {
			unit = "a"
		}
		sql = sql.Interval("?"+string(unit), filter.Interval)
	}
	if len(agg.Fill) > 0 {
		sql = sql.Fill(agg.Fill)
	}
	if filter.PartitionBy != "" {
		sql = sql.PartitionBys(partitionBy)
	}
	return sql, nil
}
func (d *DeviceDataRepo) ToPropertyData2(ctx context.Context, agg msgThing.PropertyAgg, m *schema.Model, dbs []map[string]any, retMap map[string]msgThing.PropertyData2) map[string]msgThing.PropertyData2 {
	for _, db := range dbs {
		data := msgThing.PropertyData2{
			DeviceName: cast.ToString(db["device_name"]),
			TenantCode: dataType.TenantCode(cast.ToString(db["tenant_code"])),
			ProjectID:  dataType.ProjectID(cast.ToInt64(db["project_id"])),
			AreaID:     dataType.AreaID(cast.ToInt64(db["area_id"])),
			AreaIDPath: dataType.AreaIDPath(cast.ToString(db["area_id_path"])),
		}
		for k, v := range db {
			if !strings.HasPrefix(k, "group_") {
				continue
			}
			delete(db, k)
			str := cast.ToString(v)
			if len(str) == 0 {
				continue
			}
			if data.BelongGroup == nil {
				data.BelongGroup = map[string]def.IDsInfo{}
			}
			if strings.HasSuffix(k, "_ids") {
				purpose := k[len("group_") : len(k)-len("_ids")]
				pp := data.BelongGroup[purpose]
				pp.IDs = utils.StrGenInt64Slice(str)
				data.BelongGroup[purpose] = pp
			} else if strings.HasSuffix(k, "_id_paths") {
				purpose := k[len("group_") : len(k)-len("_id_paths")]
				pp := data.BelongGroup[purpose]
				pp.IDPaths = utils.StrGenStrSlice(str)
				data.BelongGroup[purpose] = pp
			}
		}
		key := utils.MarshalNoErr(data)
		ret, ok := retMap[key]
		if !ok {
			ret = data
		}
		value := msgThing.PropertyAggData{
			Identifier: agg.DataID,
			TsWindow:   cast.ToTime(db["ts_window"]),
			Values:     map[string]msgThing.PropertyDataDetail{},
		}
		p := m.Property[agg.DataID]
		for k, v := range db {
			if strings.HasSuffix(k, "_param") {
				argFunc := k[:len(k)-len("_param")]
				vv := msgThing.PropertyDataDetail{
					Param:     v,
					TimeStamp: cast.ToTime(db[argFunc+"_ts"]),
				}
				pp, err := p.Define.FmtValue(vv.Param)
				if err != nil {
					logx.WithContext(ctx).Error("FmtValue", err)
				} else {
					vv.Param = pp
				}
				if b, ok := vv.Param.(bool); ok {
					vv.Param = cast.ToInt64(b)
				}
				if ts, ok := db[argFunc+"ts"]; ok {
					vv.TimeStamp = cast.ToTime(ts)
				}
				value.Values[argFunc] = vv
			}
		}
		ret.Values = append(ret.Values, value)
		retMap[key] = ret
	}

	return retMap
}
