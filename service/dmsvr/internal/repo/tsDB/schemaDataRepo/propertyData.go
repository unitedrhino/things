package schemaDataRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *schema.Property, productID string, deviceName string,
	property *msgThing.Param, timestamp time.Time, optional msgThing.Optional) error {
	err := d.GenInsertPropertySql(ctx, t, productID, deviceName, property, timestamp, optional)
	if err != nil {
		return err
	}
	return nil
}

func (d *DeviceDataRepo) GenInsertPropertySql(ctx context.Context, p *schema.Property, productID string, deviceName string,
	property *msgThing.Param, timestamp time.Time, optional msgThing.Optional) (err error) {
	var ars = map[string]any{}
	switch property.Define.Type {
	case schema.DataTypeArray:
		genArrSql := func(Identifier string, num int, v any) error {
			ars[schema.GenArray(Identifier, num)] = v
			if property.Define.ArrayInfo.Type == schema.DataTypeStruct {
				vv, ok := property.Value.(map[string]msgThing.Param)
				if !ok {
					return errors.Parameter.AddMsg("结构体类型传参错误")
				}
				vvv, err := msgThing.ToVal(vv)
				if err != nil {
					return err
				}
				ars[schema.GenArray(Identifier, num)] = vvv
			}
			if optional.OnlyCache || !tsDB.CheckIsChange(ctx, d.kv, devices.Core{
				ProductID: productID, DeviceName: deviceName}, p, msgThing.PropertyData{
				Identifier: schema.GenArray(Identifier, num),
				Param:      ars[schema.GenArray(Identifier, num)],
				TimeStamp:  timestamp,
			}) {
				return nil
			}

			id := GetArrayID(Identifier, num)
			pp := Property{
				ProductID:  productID,
				DeviceName: deviceName,
				Timestamp:  timestamp,
				Identifier: id,
			}
			switch property.Define.ArrayInfo.Type {
			case schema.DataTypeStruct:
				d.asyncPropertyStructArray.AsyncInsert(&PropertyStructArray{
					Property: pp,
					Param:    ars[schema.GenArray(Identifier, num)].(map[string]any),
					Pos:      int64(num),
				})
			case schema.DataTypeBool:
				d.asyncPropertyBoolArray.AsyncInsert(&PropertyBoolArray{
					Property: pp, Param: cast.ToBool(v), Pos: int64(num),
				})
			case schema.DataTypeInt:
				d.asyncPropertyIntArray.AsyncInsert(&PropertyIntArray{
					Property: pp, Param: cast.ToInt64(v), Pos: int64(num),
				})
			case schema.DataTypeString:
				d.asyncPropertyStringArray.AsyncInsert(&PropertyStringArray{
					Property: pp, Param: cast.ToString(v), Pos: int64(num),
				})
			case schema.DataTypeFloat:
				d.asyncPropertyFloatArray.AsyncInsert(&PropertyFloatArray{
					Property: pp, Param: cast.ToFloat64(v), Pos: int64(num),
				})
			case schema.DataTypeTimestamp:
				d.asyncPropertyTimestampArray.AsyncInsert(&PropertyTimestampArray{
					Property: pp, Param: cast.ToInt64(v), Pos: int64(num),
				})
			case schema.DataTypeEnum:
				d.asyncPropertyEnumArray.AsyncInsert(&PropertyEnumArray{
					Property: pp, Param: cast.ToInt64(v), Pos: int64(num),
				})
			}
			return nil
		}

		switch val := property.Value.(type) {
		case []any: //这种是数组的所有值一起上传的
			for i, v := range val {
				err := genArrSql(property.Identifier, i, v)
				if err != nil {
					return err
				}
			}
		default:
			Identifier, num, ok := schema.GetArray(property.Identifier)
			if !ok {
				return errors.Parameter.AddDetail("不是数组")
			}
			err := genArrSql(Identifier, num, val)
			if err != nil {
				return err
			}
		}
	default:
		ars[property.Identifier] = property.Value
		if property.Define.Type == schema.DataTypeStruct {
			vv, ok := property.Value.(map[string]msgThing.Param)
			if !ok {
				return errors.Parameter.AddMsg("结构体类型传参错误")
			}
			vvv, err := msgThing.ToVal(vv)
			if err != nil {
				return err
			}
			ars[property.Identifier] = vvv
		}
		if optional.OnlyCache || !tsDB.CheckIsChange(ctx, d.kv, devices.Core{
			ProductID: productID, DeviceName: deviceName}, p, msgThing.PropertyData{
			Identifier: property.Identifier,
			Param:      ars[property.Identifier],
			TimeStamp:  timestamp,
		}) {
			break
		}
		pp := Property{
			ProductID:  productID,
			DeviceName: deviceName,
			Timestamp:  timestamp,
			Identifier: property.Identifier,
		}
		switch property.Define.Type {
		case schema.DataTypeStruct:
			d.asyncPropertyStruct.AsyncInsert(&PropertyStruct{
				Property: pp,
				Param:    ars[property.Identifier].(map[string]any),
			})
		case schema.DataTypeBool:
			d.asyncPropertyBool.AsyncInsert(&PropertyBool{
				Property: pp, Param: cast.ToBool(property.Value),
			})
		case schema.DataTypeInt:
			d.asyncPropertyInt.AsyncInsert(&PropertyInt{
				Property: pp, Param: cast.ToInt64(property.Value),
			})
		case schema.DataTypeString:
			d.asyncPropertyString.AsyncInsert(&PropertyString{
				Property: pp, Param: cast.ToString(property.Value),
			})
		case schema.DataTypeFloat:
			d.asyncPropertyFloat.AsyncInsert(&PropertyFloat{
				Property: pp, Param: cast.ToFloat64(property.Value),
			})
		case schema.DataTypeTimestamp:
			d.asyncPropertyTimestamp.AsyncInsert(&PropertyTimestamp{
				Property: pp, Param: cast.ToInt64(property.Value),
			})
		case schema.DataTypeEnum:
			d.asyncPropertyEnum.AsyncInsert(&PropertyEnum{
				Property: pp, Param: cast.ToInt64(property.Value),
			})
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
			err = d.kv.Hset(tsDB.GenRedisPropertyLastKey(productID, deviceName), k, data.String())
			if err != nil {
				log.Error(err)
			}
			retStr, err := d.kv.Hget(tsDB.GenRedisPropertyFirstKey(productID, deviceName), k)
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
			err = d.kv.Hset(tsDB.GenRedisPropertyFirstKey(productID, deviceName), k, data.String())
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

func (d *DeviceDataRepo) GetLatestPropertyDataByID(ctx context.Context, p *schema.Property, filter msgThing.LatestFilter) (*msgThing.PropertyData, error) {
	retStr, err := d.kv.HgetCtx(ctx, tsDB.GenRedisPropertyLastKey(filter.ProductID, filter.DeviceName), filter.DataID)

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
			Page: def.PageInfo2{Size: 1},

			DataID: filter.DataID,
			Order:  stores.OrderDesc,
		})
	if len(dds) == 0 || err != nil {
		return nil, err
	}
	vv, er := msgThing.GetVal(&p.Define, dds[0].Param)
	if er == nil {
		dds[0].Param = vv
	}
	d.kv.HsetCtx(ctx, tsDB.GenRedisPropertyLastKey(filter.ProductID, filter.DeviceName), filter.DataID, dds[0].String())
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
		err := d.GenInsertPropertySql(ctx, p, productID, deviceName, &param, timestamp, optional)
		if err != nil {
			return errors.Database.AddDetailf(
				"DeviceDataRepo.InsertPropertiesData.InsertPropertyData identifier:%v param:%v err:%v",
				identifier, param, err)
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
		db  = d.db.WithContext(ctx).Model(getModel(p.Define))
		//needJoinAreaID bool
		//needJoinGroup  bool
	)

	if filter.ArgFunc == "" {
		if filter.Order != stores.OrderAsc {
			db = db.Order("ts desc")
		}
		if filter.PartitionBy != "" {
			var (
				selects []string
				groups  = ""
			)

			selects, groups, db = d.handlePartition(filter.PartitionBy, db)
			if len(selects) != 0 {
				db = db.Select(selects)
			}
			if groups != "" {
				db = db.Group(groups)
			}
		}
		db = db.Table(getTableName(p.Define) + " as tb")
	} else {
		db, err = d.getPropertyArgFuncSelect(ctx, db, p, filter)
		if err != nil {
			return nil, err
		}
		filter.Page.Size = 0
	}

	_, num, ok := schema.GetArray(filter.DataID)
	if ok {

		db = db.Where("pos=?", num)
	}
	db = db.Where("tb.identifier=?", filter.DataID)
	db = d.fillFilter(ctx, db, filter.Filter)
	db = filter.Page.FmtSql2(db)
	var retProperties []*msgThing.PropertyData
	var retDatabase = []map[string]any{}
	err = db.Find(&retDatabase).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	for _, v := range retDatabase {
		retProperties = append(retProperties, d.ToPropertyData(ctx, filter.NoFirstTs, filter.DataID, p, v))
	}
	return retProperties, nil
}

func (d *DeviceDataRepo) ToPropertyData(ctx context.Context, noFirstTs bool, id string, p *schema.Property, db map[string]any) *msgThing.PropertyData {
	data := msgThing.PropertyData{
		DeviceName: cast.ToString(db["device_name"]),
		Identifier: id,
		Param:      db["param"],
		TimeStamp:  cast.ToTime(db["ts_window"]),
	}
	pp, err := p.Define.FmtValue(data.Param)
	if err != nil {
		logx.WithContext(ctx).Error("FmtValue", err)
	} else {
		data.Param = pp
	}
	if b, ok := data.Param.(bool); ok {
		data.Param = cast.ToInt64(b)
	}
	if db["ts"] != nil && (noFirstTs || data.TimeStamp.IsZero()) {
		data.TimeStamp = cast.ToTime(db["ts"])
	}
	if db["tenant_code"] != nil {
		data.TenantCode = dataType.TenantCode(cast.ToString(db["tenant_code"]))
	}
	if db["project_id"] != nil {
		data.ProjectID = dataType.ProjectID(cast.ToInt64(db["project_id"]))
	}
	if db["area_id"] != nil {
		data.AreaID = dataType.AreaID(cast.ToInt64(db["area_id"]))
	}
	if db["area_id_path"] != nil {
		data.AreaIDPath = dataType.AreaIDPath(cast.ToString(db["area_id_path"]))
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
	return &data
}

/*
SELECT

	FROM_UNIXTIME(
	        UNIX_TIMESTAMP(@start) +
	        FLOOR((UNIX_TIMESTAMP(ts) - UNIX_TIMESTAMP(@start)) / @interval) * @interval
	) AS hour_window,
	device_name,
	SUM(param) AS total_param

FROM

	dm_time_model_property_float

WHERE

	ts >= @start

GROUP BY

	hour_window, device_name;
*/
func (d *DeviceDataRepo) getPropertyArgFuncSelect(
	ctx context.Context, db *stores.DB, p *schema.Property,
	filter msgThing.FilterOpt) (*stores.DB, error) {
	var start = "0000-01-01 0:00:00"
	if filter.Page.TimeStart != 0 {
		start = time.UnixMilli(filter.Page.TimeStart).Format("2006-01-02 15:04:05")
	}
	var (
		selects []string
		groups  = ""
	)
	if filter.PartitionBy != "" {
		selects, groups, db = d.handlePartition(filter.PartitionBy, db)
	}
	arg := func(ts, param string) string {
		if ts == "" {
			ts = "ts"
		}
		if param == "" {
			param = "param"
		}
		switch filter.ArgFunc {
		case "first":
			if filter.NoFirstTs {
				return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1]  AS param, (ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1] AS ts ",
					param, ts, ts, ts)
			}
			return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s ASC))[1]  AS param", param, ts)
		case "last":
			if filter.NoFirstTs {
				return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1]  AS param, (ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1] AS ts ",
					param, ts, ts, ts)
			}
			return fmt.Sprintf("(ARRAY_AGG(tb.%s ORDER BY tb.%s desc))[1]  AS param", param, ts)
		default:
			return fmt.Sprintf("%s(%s)  AS param", filter.ArgFunc, param)
		}
	}
	if filter.Interval != 0 {
		if filter.IntervalUnit == "" {
			filter.IntervalUnit = def.TimeUnitS
		}

		if stores.GetTsDBType() == conf.Pgsql && utils.SliceIn(filter.ArgFunc, "first", "last", "min", "max", "count", "sum") &&
			utils.SliceIn(filter.IntervalUnit, def.TimeUnitD, def.TimeUnitH, def.TimeUnitN, def.TimeUnitW, def.TimeUnitY) {
			//pg的 timescale走视图优化
			switch filter.IntervalUnit {
			case def.TimeUnitD, def.TimeUnitH:
				selects = append(selects, "ts as ts_window")
				db = db.Table(getTableName(p.Define) + "_" + filter.IntervalUnit.ToPgStr() + " as tb")
			default:
				selects = append(selects, fmt.Sprintf(`time_bucket('%v %s', ts)  AS ts_window`,
					filter.Interval, filter.IntervalUnit.ToPgStr()))
				db = db.Table(getTableName(p.Define) + "_day as tb")
			}
			selects = append(selects, arg(filter.ArgFunc+"_ts", filter.ArgFunc+"_param"))
			db = db.Select(selects).Group("ts_window")
		} else {
			db = db.Table(getTableName(p.Define) + " as tb")
			switch stores.GetTsDBType() {
			case conf.Pgsql:
				selects = append(selects, fmt.Sprintf(`time_bucket('%v %s', ts)  AS ts_window, %s `,
					filter.Interval, filter.IntervalUnit.ToPgStr(), arg("", "")))

				db = db.Select(selects)
			default:
				interval := int64(filter.IntervalUnit.ToDuration(filter.Interval) / time.Second)
				selects = append(selects, fmt.Sprintf(` FROM_UNIXTIME(UNIX_TIMESTAMP('%s') + FLOOR((UNIX_TIMESTAMP(ts) - UNIX_TIMESTAMP('%s')) / %v) * %v ) AS ts_window, %s AS param`,
					start, start, interval, interval, arg("", "")))

				db = db.Select(selects)
			}
			db = db.Group("ts_window")
		}

	} else {
		selects = append(selects, fmt.Sprintf("%s(param) as param", filter.ArgFunc))

		db = db.Table(getTableName(p.Define) + " as tb").Select(selects)
	}
	if groups != "" {
		db = db.Group(groups)
	}
	return db, nil
}
func (d *DeviceDataRepo) handlePartition(partitionBy string, db *stores.DB) (selects []string, groups string, db2 *stores.DB) {
	if partitionBy == "" {
		return nil, "", db
	}
	partitionBy = utils.CamelCaseToUdnderscore(partitionBy)
	ps := strings.Split(partitionBy, ",")
	var finalp []string
	for _, p := range ps {
		if strings.Contains(p, "area_id") {
			db = db.Joins("left join dm_device_info as di on tb.product_id = di.product_id and tb.device_name = di.device_name")
			selects = append(selects, "di.area_id as area_id")
			finalp = append(finalp, p)
		} else if strings.Contains(p, "group") {
			db = db.Joins("left join dm_group_device as di on tb.product_id = di.product_id and tb.device_name = di.device_name")
			selects = append(selects, "di.group_id  as group_id,di.purpose as group_purpose")

			finalp = append(finalp, "group_purpose,group_id")
		} else if p == "device_name" {
			selects = append(selects, " tb.device_name as device_name")
			finalp = append(finalp, p)
		} else {
			finalp = append(finalp, p)
		}
	}
	return selects, strings.Join(finalp, ","), db
}
func (d *DeviceDataRepo) fillFilter(ctx context.Context,
	db *stores.DB, filter msgThing.Filter) *stores.DB {
	if len(filter.DeviceNames) != 0 {
		db = db.Where("tb.device_name in ?", filter.DeviceNames)
	}
	db = tsDB.GroupFilter2(db, filter.BelongGroup)
	if len(filter.ProductIDs) != 0 {
		db = db.Where("tb.product_id IN ?", filter.ProductIDs)
	} else if filter.ProductID != "" {
		db = db.Where("tb.product_id = ?", filter.ProductID)
	}
	subQuery := d.db.WithContext(ctxs.WithDefaultAllProject(ctx)).Table("dm_device_info").Model(&relationDB.DmDeviceInfo{}).Select("product_id, device_name")
	var hasDeviceJoin bool
	if filter.ProjectID != 0 {
		subQuery = subQuery.Where("project_id=?", filter.ProjectID)
		hasDeviceJoin = true
	}
	if filter.TenantCode != "" {
		subQuery = subQuery.Where("tenant_code=?", filter.TenantCode)
		hasDeviceJoin = true
	}
	if filter.AreaID != 0 {
		subQuery = subQuery.Where("area_id=?", filter.AreaID)
		hasDeviceJoin = true
	}
	if filter.AreaIDPath != "" {
		subQuery = subQuery.Where("area_id_path like ?", filter.AreaIDPath+"%")
		hasDeviceJoin = true
	}
	if len(filter.AreaIDs) != 0 {
		subQuery = subQuery.Where("area_id in ?", filter.AreaIDs)
		hasDeviceJoin = true
	}
	if hasDeviceJoin {
		db = db.Where("(tb.product_id, tb.device_name) in (?)",
			subQuery)
	}
	return db
}

func (d *DeviceDataRepo) GetPropertyCountByID(
	ctx context.Context, p *schema.Property,
	filter msgThing.FilterOpt) (int64, error) {
	var (
		err error
		db  = d.db.WithContext(ctx).Table(getTableName(p.Define) + " as tb")
	)
	db = schema.WhereArray(db, filter.DataID, "pos")
	db = db.Where("tb.identifier=?", filter.DataID)
	db = d.fillFilter(ctx, db, filter.Filter)
	db = filter.Page.FmtSql2(db)
	var total int64
	err = db.Count(&total).Error
	return total, stores.ErrFmt(err)
}
