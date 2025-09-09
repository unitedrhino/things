package schemaDataRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) InsertPropertyData(ctx context.Context, t *schema.Property, productID string, deviceName string,
	property *msgThing.Param, timestamp time.Time, optional msgThing.Optional) error {
	sql, args, err := d.GenInsertPropertySql(ctx, t, productID, deviceName, property, timestamp, optional)
	if err != nil || sql == "" {
		return err
	}
	d.t.AsyncInsert(sql, args...)
	return nil
}

func (d *DeviceDataRepo) InsertPropertiesData(ctx context.Context, t *schema.Model, productID string, deviceName string,
	params map[string]msgThing.Param, timestamp time.Time, optional msgThing.Optional) error {
	var startTime = time.Now()
	defer func() {
		logx.WithContext(ctx).WithDuration(time.Now().Sub(startTime)).
			Infof("DeviceDataRepo.InsertPropertiesData")
	}()
	var sp = map[string]any{}
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
		if !optional.OnlyCache && sql1 != "" {
			d.t.AsyncInsert(sql1, args1...)
		}
		if len(sql1) > 0 {
			sp[identifier], _ = param.ToVal()
		}
	}
	//if len(sp) != 0 {
	//relationDB.NewShadowRepo(ctx).AsyncUpdate(ctx, shadow.NewInfo(productID, deviceName, sp, &timestamp))
	//}
	return nil
}

func (d *DeviceDataRepo) GenInsertPropertySql(ctx context.Context, p *schema.Property, productID string, deviceName string,
	property *msgThing.Param, timestamp time.Time, optional msgThing.Optional) (sql string, args []any, err error) {
	var ars = map[string]any{}

	switch property.Define.Type {
	case schema.DataTypeArray:
		genArrSql := func(Identifier string, num int, v any) error {
			ts := "`product_id` ,`device_name` ,`_data_id`,`_num`,`" + PropertyType + "`," +
				" `tenant_code` ,`project_id` ,`area_id`,`area_id_path`"
			tagKeys, tagVals := tdengine.GenTagsParams(ts, d.groupConfigs, optional.BelongGroup)

			id := GetArrayID(Identifier, num)
			k := schema.GenArray(Identifier, num)
			ars[k] = v

			// 使用缓存管理器检查是否需要更新
			if !d.cacheManager.CheckIsChange(ctx, devices.Core{ProductID: productID, DeviceName: deviceName}, p, msgThing.PropertyData{
				Identifier: k,
				Param:      v,
				TimeStamp:  timestamp,
			}) { //如果为false,则无需更新
				return nil
			}
			switch vv := v.(type) {
			case map[string]msgThing.Param, map[string]any:
				var paramPlaceholder string
				var paramIds string
				var paramValList []any
				var err error
				switch v2 := v.(type) {
				case map[string]msgThing.Param:
					paramPlaceholder, paramIds, paramValList, err = GenParams(v2)
				case map[string]any:
					paramPlaceholder, paramIds, paramValList, err = GenParams2(v2)
				}
				if err != nil {
					return err
				}
				sql += fmt.Sprintf(" %s using %s (%s)tags('%s','%s','%s',%d,'%s','%s',%d,%d,'%s' %s) (`ts`, %s) values (?,%s) ",
					d.GetPropertyTableName(productID, deviceName, id),
					d.GetPropertyStableName(p, productID, Identifier), tagKeys, productID, deviceName, Identifier, num, p.Define.Type, optional.TenantCode, optional.ProjectID,
					optional.AreaID, optional.AreaIDPath, tagVals,
					paramIds, paramPlaceholder)
				args = append([]any{timestamp}, paramValList...)
			default:
				sql += fmt.Sprintf(" %s using %s (%s)tags('%s','%s','%s',%d,'%s','%s',%d,%d,'%s' %s)(`ts`, `param`) values (?,?) ",
					d.GetPropertyTableName(productID, deviceName, id),
					d.GetPropertyStableName(p, productID, Identifier), tagKeys,
					productID, deviceName, Identifier, num, p.Define.Type, optional.TenantCode, optional.ProjectID,
					optional.AreaID, optional.AreaIDPath, tagVals)
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
		// 使用缓存管理器检查是否需要更新
		if !d.cacheManager.CheckIsChange(ctx, devices.Core{ProductID: productID, DeviceName: deviceName}, p, msgThing.PropertyData{
			Identifier: property.Identifier,
			Param:      property.Value,
			TimeStamp:  timestamp,
		}) { //如果为false,则无需更新
			break
		}
		ts := "`product_id`,`device_name`,`_data_id`,`" + PropertyType + "` ," +
			" `tenant_code`  ,`project_id` ,`area_id` ,`area_id_path` "
		tagKeys, tagVals := tdengine.GenTagsParams(ts, d.groupConfigs, optional.BelongGroup)

		switch property.Value.(type) {
		case map[string]msgThing.Param:
			paramPlaceholder, paramIds, paramValList, err := GenParams(property.Value.(map[string]msgThing.Param))
			if err != nil {
				return "", nil, err
			}
			sql = fmt.Sprintf(" %s using %s (%s)tags('%s','%s','%s','%s','%s',%d,%d,'%s' %s) (`ts`, %s) values (?,%s) ",
				d.GetPropertyTableName(productID, deviceName, property.Identifier),
				d.GetPropertyStableName(p, productID, property.Identifier), tagKeys, productID, deviceName, property.Identifier, p.Define.Type, optional.TenantCode, optional.ProjectID,
				optional.AreaID, optional.AreaIDPath, tagVals,
				paramIds, paramPlaceholder)
			args = append([]any{timestamp}, paramValList...)
			ars[property.Identifier], _ = msgThing.ToVal(property.Value.(map[string]msgThing.Param))

		default:
			var (
				param = property.Value
			)
			sql = fmt.Sprintf(" %s using %s (%s)tags('%s','%s','%s','%s','%s',%d,%d,'%s' %s)(`ts`, `param`) values (?,?) ",
				d.GetPropertyTableName(productID, deviceName, property.Identifier),
				d.GetPropertyStableName(p, productID, property.Identifier), tagKeys,
				productID, deviceName, property.Identifier, p.Define.Type, optional.TenantCode, optional.ProjectID,
				optional.AreaID, optional.AreaIDPath, tagVals)
			args = append(args, timestamp, param)
		}
	}
	f := func(ctx context.Context) {
		// 使用缓存管理器更新属性缓存
		err := d.cacheManager.UpdatePropertyCache(ctx, productID, deviceName, p, ars, timestamp)
		if err != nil {
			logx.WithContext(ctx).Errorf("更新属性缓存失败: %v", err)
		}
	}
	if !optional.Sync {
		ctxs.GoNewCtx(ctx, f)
	} else {
		f(ctx)
	}

	return
}

// GenParams 返回占位符?,?,?,? 参数id名:aa,bbb,ccc 参数值列表
func GenParams(params map[string]msgThing.Param) (string, string, []any, error) {
	if len(params) == 0 {
		//使用这个函数前必须要判断参数的个数是否大于0
		return "", "", nil, errors.Parameter.AddMsgf("SchemaDataRepo|GenParams|params num == 0")
	}
	var (
		paramPlaceholder = strings.Repeat("?,", len(params))
		paramValList     []any //参数值列表
		paramIds         []string
	)
	//将最后一个?去除
	paramPlaceholder = paramPlaceholder[:len(paramPlaceholder)-1]
	for k, vv := range params {
		v, _ := vv.ToVal()
		paramIds = append(paramIds, "`"+k+"`")
		if _, ok := v.([]any); !ok {
			paramValList = append(paramValList, v)
		} else { //如果是数组类型,需要序列化为json
			param, err := json.Marshal(v)
			if err != nil {
				return "", "", nil, errors.System.AddDetail("param json parse failure")
			}
			paramValList = append(paramValList, param)
		}
	}
	return paramPlaceholder, strings.Join(paramIds, ","), paramValList, nil
}

func GenParams2(params map[string]any) (string, string, []any, error) {
	if len(params) == 0 {
		//使用这个函数前必须要判断参数的个数是否大于0
		return "", "", nil, errors.Parameter.AddMsgf("SchemaDataRepo|GenParams|params num == 0")
	}
	var (
		paramPlaceholder = strings.Repeat("?,", len(params))
		paramValList     []any //参数值列表
		paramIds         []string
	)
	//将最后一个?去除
	paramPlaceholder = paramPlaceholder[:len(paramPlaceholder)-1]
	for k, vv := range params {
		v := vv
		paramIds = append(paramIds, "`"+k+"`")
		if _, ok := v.([]any); !ok {
			paramValList = append(paramValList, v)
		} else { //如果是数组类型,需要序列化为json
			param, err := json.Marshal(v)
			if err != nil {
				return "", "", nil, errors.System.AddDetail("param json parse failure")
			}
			paramValList = append(paramValList, param)
		}
	}
	return paramPlaceholder, strings.Join(paramIds, ","), paramValList, nil
}
