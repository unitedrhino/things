package schemaDataRepo

import (
	"encoding/json"
	"github.com/i-Things/things/shared/domain/schema"
	schema2 "github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/spf13/cast"
)

func ToEventData(db map[string]any) *schema2.EventData {
	var (
		params   map[string]any
		paramStr = cast.ToString(db["param"])
	)
	err := json.Unmarshal([]byte(paramStr), &params)
	if err != nil {
		return nil
	}
	data := schema2.EventData{
		ID:        cast.ToString(db["eventID"]),
		Type:      cast.ToString(db["eventType"]),
		Params:    params,
		TimeStamp: cast.ToTime(db["ts"]),
	}
	return &data
}

func ToPropertyData(id string, db map[string]any) *schema2.PropertyData {
	propertyType := db[PROPERTY_TYPE]
	switch propertyType {
	case string(schema.DataTypeStruct):
		data := schema2.PropertyData{
			ID:        id,
			Param:     nil,
			TimeStamp: cast.ToTime(db["ts"]),
		}
		delete(db, "ts")
		delete(db, "deviceName")
		delete(db, PROPERTY_TYPE)
		data.Param = db
		return &data
	case string(schema.DataTypeArray):
		paramStr := cast.ToString(db["param"])
		var param []any
		json.Unmarshal([]byte(paramStr), &param)
		data := schema2.PropertyData{
			ID:        id,
			Param:     param,
			TimeStamp: cast.ToTime(db["ts"]),
		}
		return &data
	default:
		data := schema2.PropertyData{
			ID:        id,
			Param:     cast.ToString(db["param"]),
			TimeStamp: cast.ToTime(db["ts"]),
		}
		return &data
	}
}
