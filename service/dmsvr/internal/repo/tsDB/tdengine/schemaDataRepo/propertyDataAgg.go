package schemaDataRepo

import (
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
)

type PropertyAggStu struct {
	NoFirstTs bool `json:"noFirstTs,optional"` //时间戳填充不填充最早的值,聚合模式使用
	Aggs      []PropertyAgg2
}
type PropertyAgg2 struct {
	P *schema.Property
	msgThing.PropertyAgg
}

//
//func (d *DeviceDataRepo) GetPropertyAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterOpt) ([]*msgThing.PropertyDatas, error) {
//
//	var (
//		err    error
//		sql    sq.SelectBuilder
//		aggMap = map[string]PropertyAggStu{} //key是表名 value是这个表下的过滤
//	)
//	if len(filter.Aggs) == 0 {
//		return nil, errors.Parameter.AddMsg("至少填写一个聚合函数")
//	}
//	for _, agg := range filter.Aggs { //暂时不考虑数组类型
//		p, ok := m.Property[agg.DataID]
//		if !ok {
//			return nil, errors.Parameter.AddMsgf("属性未定义:%v", agg.DataID)
//		}
//		tb := d.GetPropertyStableName(p, filter.ProductID, agg.DataID)
//		v := aggMap[tb+cast.ToString(agg.NoFirstTs)]
//		v.Aggs = append(v.Aggs, PropertyAgg2{
//			PropertyAgg: agg,
//			P:           p,
//		})
//	}
//	var wait sync.WaitGroup
//	for tbName, agg := range aggMap {
//		wait.Add(1)
//		utils.Go(ctx, func() {
//			defer wait.Done()
//			sql, err = d.getPropertyArgFuncSelect2(ctx, agg, filter)
//			if err != nil {
//				logx.WithContext(ctx).Errorf(err.Error())
//				return
//			}
//			sql = sql.From(tbName)
//			sql = d.fillFilter(sql, filter)
//			sql = filter.Page.FmtSql(sql)
//			sqlStr, value, err := sql.ToSql()
//			if err != nil {
//				logx.WithContext(ctx).Errorf(err.Error())
//				return
//			}
//			rows, err := d.t.QueryContext(ctx, sqlStr, value...)
//			if err != nil {
//				logx.WithContext(ctx).Errorf("sql:%v err:%v", sqlStr, err.Error())
//				return
//			}
//			var datas []map[string]any
//			stores.Scan(rows, &datas)
//		})
//	}
//	{
//		sql, err = d.getPropertyArgFuncSelect2(ctx, p, filter)
//		if err != nil {
//			return nil, err
//		}
//		filter.Page.Size = 0
//	}
//	dataID := filter.DataID
//	id, num, ok := schema.GetArray(filter.DataID)
//	if ok {
//		dataID = id
//		sql = sql.Where("`_num`=?", num)
//	}
//	sql = sql.From(d.GetPropertyStableName(p, filter.ProductID, dataID))
//	sql = d.fillFilter(sql, filter)
//	sql = filter.Page.FmtSql(sql)
//
//	sqlStr, value, err := sql.ToSql()
//	if err != nil {
//		return nil, err
//	}
//	rows, err := d.t.QueryContext(ctx, sqlStr, value...)
//	if err != nil {
//		return nil, errors.Fmt(err).AddDetailf("sql:%v", sqlStr)
//	}
//	var datas []map[string]any
//	stores.Scan(rows, &datas)
//	retProperties := make([]*msgThing.PropertyData, 0, len(datas))
//	for _, v := range datas {
//		retProperties = append(retProperties, d.ToPropertyData(ctx, filter.DataID, p, v))
//	}
//	return retProperties, err
//}
//
//func (d *DeviceDataRepo) getPropertyArgFuncSelect2(
//	ctx context.Context, agg []msgThing.PropertyAgg,
//	filter msgThing.FilterOpt) (sq.SelectBuilder, error) {
//	var (
//		sql sq.SelectBuilder
//	)
//	deviceName := ",`device_name` "
//	partitionBy := utils.CamelCaseToUdnderscore(filter.PartitionBy)
//	if !strings.Contains(partitionBy, "device_name") { //如果没有传partition by 会报错
//		deviceName = ""
//	}
//	pb := partitionBy
//	if partitionBy != "" {
//		pb = "," + pb
//	}
//	var selects []string
//	if partitionBy != "" {
//		selects = append(selects, partitionBy)
//	}
//	for _, agg := range aggs {
//		ts := "FIRST(`ts`)  AS ts "
//		if filter.Interval != 0 {
//			ts = "_wstart AS ts "
//		} else if filter.NoFirstTs {
//			ts = "`ts` "
//		}
//	}
//
//	if p.Define.Type == schema.DataTypeStruct {
//		sql = sq.Select(ts+deviceName+pb, d.GetSpecsColumnWithArgFunc(p.Define.Specs, filter.ArgFunc))
//	} else {
//		sql = sq.Select(ts+deviceName+pb, fmt.Sprintf("%s(`param`) as param", filter.ArgFunc))
//	}
//	if filter.Interval != 0 {
//		var unit = filter.IntervalUnit
//		if unit == "" {
//			unit = "a"
//		}
//		sql = sql.Interval("?"+string(unit), filter.Interval)
//	}
//	if len(filter.Fill) > 0 {
//		sql = sql.Fill(filter.Fill)
//	}
//	if filter.PartitionBy != "" {
//		sql = sql.PartitionBys(partitionBy)
//	}
//	return sql, nil
//}
