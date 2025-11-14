package schemaDataRepo

import (
	"context"
	"fmt"
	"strings"

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
)

func (d *DeviceDataRepo) GetPropertyLatestAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterLatestAggOpt) ([]*msgThing.PropertyLatestData, error) {
	var (
		err error
	)
	if len(filter.Aggs) == 0 {
		return nil, errors.Parameter.AddMsg("至少填写一个聚合函数")
	}
	for _, agg := range filter.Aggs {
		_, ok := m.Property[agg.DataID]
		if !ok {
			id, _, ok := schema.GetArray(agg.DataID)
			_, ok = m.Property[id]
			if !ok {
				return nil, errors.Parameter.AddMsg("标识符未找到")
			}
		}
	}
	var retMap = map[string]msgThing.PropertyLatestData{}
	for _, agg := range filter.Aggs {
		p, ok := m.Property[agg.DataID]
		if !ok {
			id, _, ok := schema.GetArray(agg.DataID)
			p, ok = m.Property[id]
			if !ok {
				return nil, errors.Parameter.AddMsg("标识符未找到")
			}
		}

		sql, err := d.getPropertyLatestArgFuncSelect(ctx, p, agg, filter)
		if err != nil {
			logx.WithContext(ctx).Errorf(err.Error())
			return nil, err
		}
		sqlStr, value, err := sql.ToSql()
		if err != nil {
			logx.WithContext(ctx).Errorf(err.Error())
			return nil, err
		}
		//sqlStr = fmt.Sprintf("select max(gpsTotalX) ,min(gpsTotalX),max(gpsTotalY) ,min(gpsTotalY) from (%s) ", sqlStr)
		rows, err := d.t.QueryContext(ctx, sqlStr, value...)
		if err != nil {
			logx.WithContext(ctx).Errorf("sql:%v err:%v", sqlStr, err.Error())
			return nil, err
		}
		var datas []map[string]any
		stores.Scan(rows, &datas)
		sdef := p.Define
		if sdef.Type == schema.DataTypeArray {
			sdef = *sdef.ArrayInfo
		}
		if sdef.Type == schema.DataTypeStruct { //todo 暂未支持
			dd, _ := schema.ParseDataID(agg.DataID)
			if dd != nil && dd.Column != "" {
				sdef = sdef.Spec[dd.Column].DataType
			}
		}
		retMap = d.ToPropertyLatestData(ctx, agg, &sdef, datas, retMap)
	}

	return utils.MapVToSlice2(retMap), err
}

func (d *DeviceDataRepo) getPropertyLatestArgFuncSelect(
	ctx context.Context, p *schema.Property, agg msgThing.PropertyAgg,
	filter msgThing.FilterLatestAggOpt) (sq.SelectBuilder, error) {
	var (
		sql sq.SelectBuilder
	)
	partitionBy := utils.CamelCaseToUdnderscore(filter.PartitionBy)
	var selects = []string{}
	if partitionBy != "" {
		selects = append(selects, partitionBy)
	}
	sdef := p.Define
	if sdef.Type == schema.DataTypeArray {
		sdef = *sdef.ArrayInfo
	}
	if sdef.Type == schema.DataTypeStruct {
		dd, _ := schema.ParseDataID(agg.DataID)
		if dd != nil && dd.Column != "" {
			selects = append(selects, fmt.Sprintf(" last(`%s`) as `%s`", dd.Column, dd.Column))
		} else {
			for _, v := range sdef.Specs {
				selects = append(selects, fmt.Sprintf(" last(`%s`) as `%s`", v.Identifier, v.Identifier))
			}
		}
	} else {
		selects = append(selects, " last(`param`) as `param`")
	}
	sql = sq.Select(selects...)
	if filter.PartitionBy != "" {
		sql = sql.PartitionBys(partitionBy)
	}
	id, _, _ := schema.GetArray(agg.DataID)
	sql = schema.WhereArray2(sql, agg.DataID, "`_num`")
	sql = sql.Where("`_data_id`=?", id)
	sql = sql.From(d.GetPropertyStableName(p, filter.ProductID, id))
	sql = d.fillFilter(sql, filter.Filter)
	//sqlStr,args,err:=sql.ToSql()
	//if err != nil {
	//	return sql, err
	//}

	var sql2 sq.SelectBuilder
	var selects2 = []string{}

	getOnCol := func(col string) {
		for _, argFunc := range agg.ArgFuncs {
			selects2 = append(selects2, fmt.Sprintf(" %s(`%s`) as `%s_param` ", argFunc, col, argFunc))
		}
	}

	if sdef.Type == schema.DataTypeStruct {
		dd, _ := schema.ParseDataID(agg.DataID)
		if dd != nil && dd.Column != "" {
			getOnCol(dd.Column)
		} else {
			selects2 = append(selects2, d.GetSpecsColumnWithArgFunc2(sdef.Specs, agg))
		}
	} else {
		getOnCol("`param`")
	}

	sql2 = sq.Select(selects2...)
	sql2 = sql2.FromSelect(sql, "tb")

	return sql2, nil
}

func (d *DeviceDataRepo) ToPropertyLatestData(ctx context.Context, agg msgThing.PropertyAgg, p *schema.Define, dbs []map[string]any, retMap map[string]msgThing.PropertyLatestData) map[string]msgThing.PropertyLatestData {
	dd, _ := schema.ParseDataID(agg.DataID)
	for _, db := range dbs {
		data := msgThing.PropertyLatestData{
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
		value := msgThing.PropertyLatestAggData{
			Identifier: agg.DataID,
			Values:     map[string]any{},
		}
		if p.Type == schema.DataTypeStruct && dd != nil && dd.Column == "" { //获取结构体所有字段
			var argCache = map[string]*structH{}
			for k, v := range db {
				if strings.HasSuffix(k, "_param") { // dataID_argFunc_param
					keys := strings.Split(k, "_")
					if len(keys) < 3 {
						continue
					}
					argFunc := keys[len(keys)-2]
					if argCache[argFunc] == nil {
						argCache[argFunc] = &structH{Values: map[string]any{}}
					}
					dataID := strings.Join(keys[:len(keys)-2], "_")
					sp := p.Spec[dataID]
					if sp == nil {
						for pk, spp := range p.Spec {
							if strings.ToLower(pk) == dataID {
								sp = spp
								dataID = pk
							}
						}
						if sp == nil {
							continue
						}
					}
					pp, err := sp.DataType.FmtValue(v)
					if err != nil {
						logx.WithContext(ctx).Error("FmtValue", err)
					} else {
						v = pp
					}
					if b, ok := v.(bool); ok {
						v = cast.ToInt64(b)
					}
					if ts, ok := db[argFunc+"ts"]; ok {
						v = cast.ToTime(ts)
					}
					argCache[argFunc].Values[dataID] = v
				}
			}
			for k, v := range argCache {
				value.Values[k] = v.Values
			}
		} else {
			for k, v := range db {
				if strings.HasSuffix(k, "_param") {
					argFunc := k[:len(k)-len("_param")]
					vv := msgThing.PropertyLogDataDetail{
						Param:     v,
						TimeStamp: cast.ToTime(db[argFunc+"_ts"]),
					}
					pp, err := p.FmtValue(vv.Param)
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
					value.Values[argFunc] = vv.Param
				}
			}
		}

		ret.Values = append(ret.Values, value)
		retMap[key] = ret
	}

	return retMap
}
