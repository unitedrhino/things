package schemaDataRepo

import (
	"encoding/json"
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
	switch propertyType {
	case string(schema.DataTypeStruct):
		data := msgThing.PropertyData{
			DeviceName: cast.ToString(db["device_name"]),
			Identifier: id,
			Param:      nil,
			TimeStamp:  cast.ToTime(db["ts"]),
		}
		delete(db, "ts")
		delete(db, "device_name")
		delete(db, PropertyType)
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
			delete(db, "ts")
			delete(db, "device_name")
			delete(db, PropertyType)
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
			return &data
		}
	default:
		data := msgThing.PropertyData{
			Identifier: id,
			DeviceName: cast.ToString(db["device_name"]),
			Param:      cast.ToString(utils.BoolToInt(db["param"])),
			TimeStamp:  cast.ToTime(db["ts"]),
		}
		return &data
	}
}
