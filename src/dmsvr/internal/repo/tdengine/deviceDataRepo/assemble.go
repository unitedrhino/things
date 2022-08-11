package deviceDataRepo

import (
	"encoding/json"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"
	"github.com/spf13/cast"
)

func ToEventData(id string, db map[string]any) *deviceData.EventData {
	var (
		params   map[string]any
		paramStr = cast.ToString(db["param"])
	)
	err := json.Unmarshal([]byte(paramStr), &params)
	if err != nil {
		return nil
	}
	data := deviceData.EventData{
		ID:        id,
		Type:      cast.ToString(db["event_type"]),
		Params:    params,
		TimeStamp: cast.ToTime(db["ts"]),
	}
	return &data
}

func ToPropertyData(id string, db map[string]any) *deviceData.PropertyData {
	propertyType := db[PROPERTY_TYPE]
	switch propertyType {
	case string(schema.STRUCT):
		data := deviceData.PropertyData{
			ID:        id,
			Param:     nil,
			TimeStamp: cast.ToTime(db["ts"]),
		}
		delete(db, "ts")
		delete(db, "device_name")
		delete(db, PROPERTY_TYPE)
		data.Param = db
		return &data
	case string(schema.ARRAY):
		paramStr := cast.ToString(db["param"])
		var param []any
		json.Unmarshal([]byte(paramStr), &param)
		data := deviceData.PropertyData{
			ID:        id,
			Param:     param,
			TimeStamp: cast.ToTime(db["ts"]),
		}
		return &data
	default:
		data := deviceData.PropertyData{
			ID:        id,
			Param:     cast.ToString(db["param"]),
			TimeStamp: cast.ToTime(db["ts"]),
		}
		return &data
	}
}
