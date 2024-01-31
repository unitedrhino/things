package schemaDataRepo

import (
	"encoding/json"
	"gitee.com/i-Things/core/shared/domain/schema"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgThing"
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
		Params:     params,
		TimeStamp:  cast.ToTime(db["ts"]),
	}
	return &data
}

func ToPropertyData(id string, db map[string]any) *msgThing.PropertyData {
	propertyType := db[PropertyType]
	switch propertyType {
	case string(schema.DataTypeStruct):
		data := msgThing.PropertyData{
			Identifier: id,
			Param:      nil,
			TimeStamp:  cast.ToTime(db["ts"]),
		}
		delete(db, "ts")
		delete(db, "device_name")
		delete(db, PropertyType)
		data.Param = db
		return &data
	case string(schema.DataTypeArray):
		paramStr := cast.ToString(db["param"])
		var param []any
		json.Unmarshal([]byte(paramStr), &param)
		data := msgThing.PropertyData{
			Identifier: id,
			Param:      param,
			TimeStamp:  cast.ToTime(db["ts"]),
		}
		return &data
	default:
		data := msgThing.PropertyData{
			Identifier: id,
			Param:      cast.ToString(db["param"]),
			TimeStamp:  cast.ToTime(db["ts"]),
		}
		return &data
	}
}
