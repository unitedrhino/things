package schemaDataRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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
			if optional.OnlyCache {
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
				vv, ok := property.Value.(map[string]msgThing.Param)
				if !ok {
					return errors.Parameter.AddMsg("结构体类型传参错误")
				}
				vvv, err := msgThing.ToVal(vv)
				if err != nil {
					return err
				}
				ars[property.Identifier] = vvv
				d.asyncPropertyStructArray.AsyncInsert(&PropertyStructArray{
					Property: pp,
					Param:    vvv,
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
		if optional.OnlyCache {
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
			vv, ok := property.Value.(map[string]msgThing.Param)
			if !ok {
				return errors.Parameter.AddMsg("结构体类型传参错误")
			}
			vvv, err := msgThing.ToVal(vv)
			if err != nil {
				return err
			}
			ars[property.Identifier] = vvv
			d.asyncPropertyStruct.AsyncInsert(&PropertyStruct{
				Property: pp,
				Param:    vvv,
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
	)
	if filter.PartitionBy != "" {
		db = db.Group(filter.PartitionBy)
	}
	if filter.ArgFunc == "" {
		if filter.Order != stores.OrderAsc {
			db = db.Order("ts desc")
		}
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
	db = d.fillFilter(db, filter)
	db = filter.Page.FmtSql2(db)
	var retProperties []*msgThing.PropertyData
	if p.Define.Type == schema.DataTypeArray {
		p.Define.Type = p.Define.ArrayInfo.Type
	}
	switch p.Define.Type {
	case schema.DataTypeBool:
		var ret = []*PropertyBool{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	case schema.DataTypeInt:
		var ret = []*PropertyInt{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	case schema.DataTypeString:
		var ret = []*PropertyString{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	case schema.DataTypeStruct:
		var ret = []*PropertyStruct{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	case schema.DataTypeFloat:
		var ret = []*PropertyFloat{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	case schema.DataTypeTimestamp:
		var ret = []*PropertyTimestamp{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	case schema.DataTypeEnum:
		var ret = []*PropertyEnum{}
		err = db.Find(&ret).Error
		if err != nil {
			return nil, stores.ErrFmt(err)
		}
		retProperties = utils.CopySlice[msgThing.PropertyData](ret)
	}
	return retProperties, nil
}

func (d *DeviceDataRepo) getPropertyArgFuncSelect(
	ctx context.Context, db *stores.DB, p *schema.Property,
	filter msgThing.FilterOpt) (*stores.DB, error) {

	if filter.Interval != 0 {
		db = db.Select("ts", fmt.Sprintf(
			" DATE_FORMAT(FROM_UNIXTIME(FLOOR(UNIX_TIMESTAMP(ts) / %v) * %v), '%Y-%m-%d %H:%i:00') AS time_interval"), fmt.Sprintf("%s(`param`) as param", filter.ArgFunc))
		db = db.Group("time_interval")
	} else {
		db = db.Select("ts", fmt.Sprintf("%s(`param`) as param", filter.ArgFunc))
	}
	return db, nil
}

func (d *DeviceDataRepo) fillFilter(
	db *stores.DB, filter msgThing.FilterOpt) *stores.DB {
	db = db.Where("product_id=?", filter.ProductID)
	if len(filter.DeviceNames) != 0 {
		db = db.Where(fmt.Sprintf("device_name in (%v)", stores.ArrayToSql(filter.DeviceNames)))
	}
	return db
}

func (d *DeviceDataRepo) GetPropertyCountByID(
	ctx context.Context, p *schema.Property,
	filter msgThing.FilterOpt) (int64, error) {
	var (
		err error
		db  = d.db.WithContext(ctx).Table(getTableName(p.Define))
	)
	_, num, ok := schema.GetArray(filter.DataID)
	if ok {
		db = db.Where("pos=?", num)
	}
	db = d.fillFilter(db, filter)
	db = filter.Page.FmtSql2(db)
	var total int64
	err = db.Count(&total).Error
	return total, stores.ErrFmt(err)
}
