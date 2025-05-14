package schemaDataRepo

import (
	"encoding/json"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"github.com/spf13/cast"
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

func ToPropertyData(id string, p *schema.Property, db map[string]any) *msgThing.PropertyData {
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
		if db["group_ids"] != nil {
			data.GroupIDs = utils.StrGenInt64Slice(cast.ToString(db["group_ids"]))
		}
		if db["group_id_paths"] != nil {
			data.GroupIDPaths = utils.StrGenStrSlice(cast.ToString(db["group_id_paths"]))
		}
		delete(db, "ts")
		delete(db, "device_name")
		delete(db, "tenant_code")
		delete(db, "project_id")
		delete(db, "area_id")
		delete(db, "area_id_path")
		delete(db, "group_ids")
		delete(db, "group_id_paths")
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
