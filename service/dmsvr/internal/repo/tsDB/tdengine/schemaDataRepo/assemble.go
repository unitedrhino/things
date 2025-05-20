package schemaDataRepo

import (
	"encoding/json"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
	"strings"
)

func ToEventData(db map[string]any) *msgThing.EventData {
	var (
		params   map[string]any
		paramStr = cast.ToString(db["param"])
	)
	err := json.Unmarshal([]byte(paramStr), &params)
	if err != nil {
		return nil
	}
	data := msgThing.EventData{
		Identifier: cast.ToString(db["event_id"]),
		Type:       cast.ToString(db["event_type"]),
		DeviceName: cast.ToString(db["device_name"]),
		Params:     params,
		TimeStamp:  cast.ToTime(db["ts"]),
	}
	return &data
}

func (d *DeviceDataRepo) ToPropertyData(id string, p *schema.Property, db map[string]any) *msgThing.PropertyData {
	propertyType := db[PropertyType]
	fill := func(data *msgThing.PropertyData) {
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
			data.TenantCode = dataType.TenantCode(cast.ToString(db["area_id_path"]))
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
		delete(db, "ts")
		delete(db, "device_name")
		delete(db, "tenant_code")
		delete(db, "project_id")
		delete(db, "area_id")
		delete(db, "area_id_path")
		delete(db, PropertyType)
	}
	switch propertyType {
	case string(schema.DataTypeStruct):
		data := msgThing.PropertyData{
			DeviceName: cast.ToString(db["device_name"]),
			Identifier: id,
			Param:      nil,
			TimeStamp:  cast.ToTime(db["ts"]),
		}

		fill(&data)
		delete(db, "_num")
		for k, v := range db {
			db[k] = utils.BoolToInt(v)
		}
		data.Param = db
		return &data
	case string(schema.DataTypeArray):
		switch p.Define.ArrayInfo.Type {
		case schema.DataTypeStruct:
			data := msgThing.PropertyData{
				Identifier: id,
				DeviceName: cast.ToString(db["device_name"]),
				TimeStamp:  cast.ToTime(db["ts"]),
			}
			fill(&data)
			delete(db, "_num")
			for k, v := range db {
				db[k] = utils.BoolToInt(v)
			}
			data.Param = db
			return &data
		default:
			data := msgThing.PropertyData{
				Identifier: id,
				DeviceName: cast.ToString(db["device_name"]),
				Param:      cast.ToString(utils.BoolToInt(db["param"])),
				TimeStamp:  cast.ToTime(db["ts"]),
			}
			fill(&data)
			return &data
		}
	default:
		data := msgThing.PropertyData{
			Identifier: id,
			DeviceName: cast.ToString(db["device_name"]),
			Param:      cast.ToString(utils.BoolToInt(db["param"])),
			TimeStamp:  cast.ToTime(db["ts"]),
		}
		fill(&data)
		return &data
	}
}
