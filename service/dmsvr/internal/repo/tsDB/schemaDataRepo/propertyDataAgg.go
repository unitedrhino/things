package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"
	"time"
)

func (d *DeviceDataRepo) GetPropertyAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterAggOpt) ([]*msgThing.PropertyData2, error) {

	var (
		err error
	)
	if len(filter.Aggs) == 0 {
		return nil, errors.Parameter.AddMsg("至少填写一个聚合函数")
	}
	for _, agg := range filter.Aggs { //暂时不考虑数组类型
		_, ok := m.Property[agg.DataID]
		if !ok {
			return nil, errors.Parameter.AddMsgf("属性未定义:%v", agg.DataID)
		}
	}
	var retMap = map[string]msgThing.PropertyData2{}
	var mutex sync.Mutex
	var wait sync.WaitGroup
	for _, agg := range filter.Aggs {
		wait.Add(1)
		utils.Go(ctx, func() {
			defer wait.Done()
			var (
				err  error
				p, _ = m.Property[agg.DataID]
				db   = d.db.WithContext(ctx).Model(getModel(p.Define))
				//needJoinAreaID bool
				//needJoinGroup  bool
			)
			db, err = d.getPropertyArgFuncSelect2(ctx, db, p, agg, filter)
			if err != nil {
				logx.WithContext(ctx).Errorf(err.Error())
				return
			}
			_, num, ok := schema.GetArray(agg.DataID)
			if ok {
				db = db.Where("pos=?", num)
			}
			if filter.TimeStart > 0 {
				db = db.Where("ts>=?", time.UnixMilli(filter.TimeStart))
			}
			if filter.TimeEnd > 0 {
				db = db.Where("ts<=?", time.UnixMilli(filter.TimeEnd))
			}
			db = db.Where("tb.identifier=?", agg.DataID)
			db = d.fillFilter(ctx, db, filter.Filter)
			//db = filter.Page.FmtSql2(db)
			var retDatabase = []map[string]any{}
			err = db.Find(&retDatabase).Error
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			retMap = d.ToPropertyData2(ctx, agg, m, retDatabase, retMap)
		})
	}
	wait.Wait()
	return utils.MapVToSlice2(retMap), err
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect2(
	ctx context.Context, db *stores.DB, p *schema.Property, agg msgThing.PropertyAgg,
	filter msgThing.FilterAggOpt) (*stores.DB, error) {
	//var start = "0000-01-01 0:00:00"
	//if filter.TimeStart != 0 {
	//	start = time.UnixMilli(filter.TimeStart).Format("2006-01-02 15:04:05")
	//}
	var (
		selects []string
		groups  = ""
	)
	if filter.PartitionBy != "" {
		selects, groups, db = d.handlePartition(filter.PartitionBy, db)
	}
	arg := func(ts, param string, argFunc string) string {
		if ts == "" {
			ts = "ts"
		}
		if param == "" {
			param = "param"
		}
		switch argFunc {
		case "first":
			if agg.NoFirstTs {
				return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1]  AS %s_param, (ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1] AS %s_ts ",
					param, ts, argFunc, ts, ts, argFunc)
			}
			return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1]  AS %s_param", param, ts, argFunc)
		case "last":
			if agg.NoFirstTs {
				return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1]  AS %s_param, (ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1] AS %s_ts ",
					param, ts, argFunc, ts, ts, argFunc)
			}
			return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1]  AS %s_param", param, ts, argFunc)
		default:
			return fmt.Sprintf("%s(%s)  AS %s_param", argFunc, param, argFunc)
		}
	}
	if filter.Interval != 0 {
		if filter.IntervalUnit == "" {
			filter.IntervalUnit = def.TimeUnitS
		}
		if stores.GetTsDBType() != conf.Pgsql { //todo 待实现

		} else {
			if utils.SliceIn(filter.IntervalUnit, def.TimeUnitD, def.TimeUnitH, def.TimeUnitN, def.TimeUnitW, def.TimeUnitY) {
				var tbName = getTableName(p.Define) + "_day as tb"
				switch filter.IntervalUnit {
				case def.TimeUnitD, def.TimeUnitH:
					if filter.IntervalUnit == def.TimeUnitH {
						tbName = getTableName(p.Define) + "_hour as tb"
					}
					selects = append(selects, "ts as ts_window")
					db = db.Where("")
					for _, argFunc := range agg.ArgFuncs {
						//pg的 timescale走视图优化
						if agg.NoFirstTs && utils.SliceIn(argFunc, "first", "last", "min", "max") {
							selects = append(selects, fmt.Sprintf(` %s_ts , %s_param `, argFunc, argFunc))
						} else {
							selects = append(selects, fmt.Sprintf(` %s_param `, argFunc))
						}
						groups = ""
					}
				default:
					db = db.Group("ts_window")
					selects = append(selects, fmt.Sprintf(`time_bucket('%v %s', ts)  AS ts_window `,
						filter.Interval, filter.IntervalUnit.ToPgStr()))
					for _, argFunc := range agg.ArgFuncs {
						selects = append(selects, arg(argFunc+"_ts", argFunc+"_param", argFunc))
					}
				}
				db = db.Select(selects)
				if groups != "" {
					db = db.Group(groups)
				}
				db = db.Table(tbName)
				return db, nil
			} else { //todo 待实现
				db = db.Table(getTableName(p.Define) + " as tb")
				selects = append(selects, fmt.Sprintf(`time_bucket('%v %s', ts)  AS ts_window, %s `,
					filter.Interval, filter.IntervalUnit.ToPgStr(), arg("", "", "argFunc")))
				db = db.Select(selects)
				db = db.Group("ts_window")
			}
		}

	} else {
		//todo 待实现
		//db = db.Table(getTableName(p.Define) + " as tb").Select(selects + fmt.Sprintf("%s(param) as param", filter.ArgFunc))
	}
	if groups != "" {
		db = db.Group(groups)
	}
	return db, nil
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
		if db["group_purpose"] != nil {
			groupPurpose := cast.ToString(db["group_purpose"])
			var idsInfo def.IDsInfo
			if db["group_id"] != nil {
				idsInfo.IDs = append(idsInfo.IDs, cast.ToInt64(db["group_id"]))
			}
			if db["group_id_path"] != nil {
				idsInfo.IDPaths = append(idsInfo.IDPaths, cast.ToString(db["group_id_path"]))
			}
			data.BelongGroup = map[string]def.IDsInfo{
				groupPurpose: idsInfo,
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
