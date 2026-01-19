package schemaDataRepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) GetLatestAllPropertyData(ctx context.Context, m *schema.Model, productID, deviceName string) ([]*msgThing.PropertyLogData, error) {
	// 使用缓存管理器获取设备所有属性的最后记录
	cacheData, err := d.cacheManager.GetPropertyAllLastRecord(ctx, m, productID, deviceName)
	if err != nil {
		logx.WithContext(ctx).Error(err)
	} else if len(cacheData) > 0 {
		// 缓存中有数据，直接返回
		return cacheData, nil
	}

	// 缓存中没有数据，从数据库加载所有属性
	var result []*msgThing.PropertyLogData
	for _, property := range m.Properties {
		// 对每个属性从数据库获取最新数据
		filter := msgThing.LatestFilter{
			ProductID:  productID,
			DeviceName: deviceName,
			DataID:     property.Identifier,
		}
		propertyData, err := d.GetLatestPropertyDataByID(ctx, &property, filter)
		if err != nil {
			logx.WithContext(ctx).Errorf("获取属性 %s 最新数据失败: %v", property.Identifier, err)
			continue
		}
		if propertyData != nil {
			result = append(result, propertyData)
		}
	}

	// 更新缓存（异步执行）
	if len(result) > 0 {
		// 将整个缓存更新逻辑放在一个异步任务中执行
		utils.Go(ctx, func() {
			// 组织数据格式以便更新缓存
			for _, data := range result {
				// 查找对应的属性定义
				var property *schema.Property
				for i := range m.Properties {
					if m.Properties[i].Identifier == data.Identifier {
						property = &m.Properties[i]
						break
					}
				}
				if property != nil {
					// 更新单个属性缓存
					err := d.cacheManager.UpdatePropertyCache(ctx, productID, deviceName, property, map[string]any{data.Identifier: data.Param}, data.TimeStamp)
					if err != nil {
						logx.WithContext(ctx).Errorf("更新属性 %s 缓存失败: %v", data.Identifier, err)
					}
				}
			}
		})
	}

	return result, nil
}

func (d *DeviceDataRepo) GetLatestPropertyDataByID(ctx context.Context, p *schema.Property, filter msgThing.LatestFilter) (*msgThing.PropertyLogData, error) {
	// 使用缓存管理器获取最后记录
	ret, err := d.cacheManager.GetPropertyLastRecord(ctx, filter.ProductID, filter.DeviceName, filter.DataID)
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	if ret != nil && ret.TimeStamp.After(time.Now().Add(-time.Hour*24)) { //只保留一个小时
		vv, er := msgThing.GetVal(&p.Define, ret.Param)
		if er == nil {
			ret.Param = vv
		}
		return ret, nil
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
	// 更新缓存
	err = d.cacheManager.UpdatePropertyCache(ctx, filter.ProductID, filter.DeviceName, p, map[string]any{filter.DataID: dds[0].Param}, dds[0].TimeStamp)
	if err != nil {
		logx.WithContext(ctx).Errorf("更新属性缓存失败: %v", err)
	}
	return dds[0], nil

}

//func (d *DeviceDataRepo) GetPropertyLogAgg(ctx context.Context, m *schema.Model, filter msgThing.FilterOpt) ([]*msgThing.PropertyDatas, error) {
//	//TODO implement me
//	panic("implement me")
//}

func (d *DeviceDataRepo) GetPropertyDataByID(
	ctx context.Context, p *schema.Property,
	filter msgThing.FilterOpt) ([]*msgThing.PropertyLogData, error) {
	if err := filter.Check(); err != nil {
		return nil, err
	}

	var (
		err error
		sql sq.SelectBuilder
	)

	if filter.ArgFunc == "" {
		h := func() bool {
			sdef := p.Define
			if sdef.Type == schema.DataTypeArray {
				sdef = *sdef.ArrayInfo
			}
			if sdef.Type == schema.DataTypeStruct {
				dd, _ := schema.ParseDataID(filter.DataID)
				if dd != nil && dd.Column != "" {
					sql = sq.Select("`ts`,`device_name`", fmt.Sprintf("`%s` as param", dd.Column))
					return true
				}
			}
			return false
		}()
		if !h {
			sql = sq.Select("*")
			if filter.Order != stores.OrderAsc {
				sql = sql.OrderBy("`ts` desc")
			}
		}
	} else {
		sql, err = d.getPropertyArgFuncSelect(ctx, p, filter)
		if err != nil {
			return nil, err
		}
		filter.Page.Size = 0
	}
	sql = schema.WhereArray2(sql, filter.DataID, "`_num`")
	id, _, _ := schema.GetArray(filter.DataID)
	sql = sql.Where("`_data_id`=?", id)
	sql = sql.From(d.GetPropertyStableName(p, filter.ProductID, id))
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
	retProperties := make([]*msgThing.PropertyLogData, 0, len(datas))
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
	partitionBy := utils.CamelCaseToUdnderscore(filter.PartitionBy)

	// 定义 selects 数组
	var selects []string

	// 添加 partition by 字段，如果没有则默认使用 device_name 避免聚合错误
	if partitionBy != "" {
		selects = append(selects, partitionBy)
	} else {
		// 如果没有指定 partitionBy，默认添加 device_name
		selects = append(selects, "`device_name`")
	}

	// 定义 getOnCol 函数用于添加聚合列
	getOnCol := func(col string) {
		//pg的 timescale走视图优化
		if filter.NoFirstTs && utils.SliceIn(filter.ArgFunc, "first", "last", "min", "max") {
			selects = append(selects, fmt.Sprintf("%s(`%s`) as param, cols(%s(`%s`),ts) as ts", filter.ArgFunc, col, filter.ArgFunc, col))
		} else {
			// 根据不同条件添加时间戳
			if filter.Interval != 0 {
				selects = append(selects, "_wstart AS ts")
			} else if filter.NoFirstTs {
				selects = append(selects, "`ts`")
			} else {
				selects = append(selects, "FIRST(`ts`) AS ts")
			}
			// 添加聚合列
			selects = append(selects, fmt.Sprintf("%s(`%s`) as param", filter.ArgFunc, col))
		}
	}

	// 处理不同的数据类型
	sdef := p.Define
	if sdef.Type == schema.DataTypeArray {
		sdef = *sdef.ArrayInfo
	}
	if sdef.Type == schema.DataTypeStruct {
		dd, _ := schema.ParseDataID(filter.DataID)
		if dd != nil && dd.Column != "" {
			getOnCol(dd.Column)
		} else {
			// 对于 struct 的多列情况，需要单独添加时间戳
			//pg的 timescale走视图优化
			if filter.NoFirstTs && utils.SliceIn(filter.ArgFunc, "first", "last", "min", "max") {
				// struct 多列暂时不支持 cols 优化，保持原有逻辑
				selects = append(selects, "`ts`")
				selects = append(selects, d.GetSpecsColumnWithArgFunc(sdef.Specs, filter.ArgFunc))
			} else {
				if filter.Interval != 0 {
					selects = append(selects, "_wstart AS ts")
				} else if filter.NoFirstTs {
					selects = append(selects, "`ts`")
				} else {
					selects = append(selects, "FIRST(`ts`) AS ts")
				}
				selects = append(selects, d.GetSpecsColumnWithArgFunc(sdef.Specs, filter.ArgFunc))
			}
		}
	} else {
		getOnCol("param")
	}

	// 使用 selects 数组构建 SQL
	sql = sq.Select(selects...)

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
	sqlData = sqlData.Where("`_data_id`=?", id)
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
