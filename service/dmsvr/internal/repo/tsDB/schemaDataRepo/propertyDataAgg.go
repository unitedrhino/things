package schemaDataRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
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
	//var wait sync.WaitGroup
	for _, agg := range filter.Aggs {
		//	wait.Add(1)
		//	utils.Go(ctx, func() {
		//		defer wait.Done()
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
			return nil, err
		}
		_, num, ok := schema.GetArray(agg.DataID)
		if ok {

			db = db.Where("pos=?", num)
		}
		//db = d.fillFilter(ctx, db, filter)
		//db = filter.Page.FmtSql2(db)
		//var retProperties []*msgThing.PropertyData
		var retDatabase = []map[string]any{}
		err = db.Find(&retDatabase).Error
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return nil, err
		}
		//for _, v := range retDatabase {
		//	retProperties = append(retProperties, d.ToPropertyData(ctx, filter.NoFirstTs, filter.DataID, p, v))
		//}
		//	})
	}
	//wait.Wait()
	return nil, err
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect2(
	ctx context.Context, db *stores.DB, p *schema.Property, agg msgThing.PropertyAgg,
	filter msgThing.FilterAggOpt) (*stores.DB, error) {
	var start = "0000-01-01 0:00:00"
	if filter.TimeStart != 0 {
		start = time.UnixMilli(filter.TimeStart).Format("2006-01-02 15:04:05")
	}
	var (
		selects = ""
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
				return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1]  AS param, (ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1] AS ts ",
					param, ts, ts, ts)
			}
			return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1]  AS param", param, ts)
		case "last":
			if agg.NoFirstTs {
				return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1]  AS param, (ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1] AS ts ",
					param, ts, ts, ts)
			}
			return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1]  AS param", param, ts)
		default:
			return fmt.Sprintf("%s(%s)  AS param", argFunc, param)
		}
	}
	if filter.Interval != 0 {
		if filter.IntervalUnit == "" {
			filter.IntervalUnit = def.TimeUnitS
		}
		if stores.GetTsDBType() == conf.Pgsql &&
			utils.SliceIn(filter.IntervalUnit, def.TimeUnitD, def.TimeUnitH, def.TimeUnitN, def.TimeUnitW, def.TimeUnitY) {
			//pg的 timescale走视图优化
			switch filter.IntervalUnit {
			case def.TimeUnitD, def.TimeUnitH:
				selects += "ts as ts_window,"
			default:
				selects += fmt.Sprintf("time_bucket('%v %s', ts)  AS ts_window,", filter.Interval, filter.IntervalUnit.ToPgStr())
			}
		}
		for _, argFunc := range agg.ArgFuncs {
			if stores.GetTsDBType() == conf.Pgsql && utils.SliceIn(argFunc, "first", "last", "min", "max", "count", "sum", "avg") &&
				utils.SliceIn(filter.IntervalUnit, def.TimeUnitD, def.TimeUnitH, def.TimeUnitN, def.TimeUnitW, def.TimeUnitY) {
				//pg的 timescale走视图优化
				switch filter.IntervalUnit {
				case def.TimeUnitD, def.TimeUnitH:
					db = db.Table(getTableName(p.Define) + "_" + filter.IntervalUnit.ToPgStr() + " as tb")
					if agg.NoFirstTs && utils.SliceIn(argFunc, "first", "last", "min", "max") {
						db = db.Select(selects + fmt.Sprintf(` %s_ts as ts, %s_param as param `, argFunc, argFunc))
					} else {
						db = db.Select(selects + fmt.Sprintf(` %s_param as param `, argFunc))
					}
					groups = ""
				default:
					db = db.Select(selects + arg(argFunc+"_ts", argFunc+"_param", argFunc))
					db = db.Table(getTableName(p.Define) + "_day as tb").Group("ts_window")
				}
			} else {
				db = db.Table(getTableName(p.Define) + " as tb")
				switch stores.GetTsDBType() {
				case conf.Pgsql:
					db = db.Select(selects + fmt.Sprintf(`time_bucket('%v %s', ts)  AS ts_window, %s `,
						filter.Interval, filter.IntervalUnit.ToPgStr(), arg("", "", argFunc)))
				default:
					interval := int64(filter.IntervalUnit.ToDuration(filter.Interval) / time.Second)
					db = db.Select(selects + fmt.Sprintf(` FROM_UNIXTIME(UNIX_TIMESTAMP('%s') + FLOOR((UNIX_TIMESTAMP(ts) - UNIX_TIMESTAMP('%s')) / %v) * %v ) AS ts_window, %s AS param`,
						start, start, interval, interval, arg("", "", argFunc)))
				}
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
