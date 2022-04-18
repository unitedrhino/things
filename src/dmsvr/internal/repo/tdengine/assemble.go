package tdengine

import (
	"encoding/json"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"time"
)

func ToEventData(id string, db map[string]interface{}) *deviceTemplate.EventData {
	var (
		params   map[string]interface{}
		paramStr = db["param"].(string)
	)
	err := json.Unmarshal([]byte(paramStr), &params)
	if err != nil {
		return nil
	}
	data := deviceTemplate.EventData{
		ID:        id,
		Type:      db["event_type"].(string),
		Params:    params,
		TimeStamp: db["ts"].(time.Time),
	}
	return &data
}

func ToPropertyData(id string, db map[string]interface{}) *deviceTemplate.PropertyData {
	propertyType := db[PROPERTY_TYPE]
	switch propertyType {
	case string(deviceTemplate.STRUCT):
		data := deviceTemplate.PropertyData{
			ID:        id,
			Param:     nil,
			TimeStamp: db["ts"].(time.Time),
		}
		delete(db, "ts")
		delete(db, "device_name")
		delete(db, PROPERTY_TYPE)
		data.Param = db
		return &data
	default:
		data := deviceTemplate.PropertyData{
			ID:        id,
			Param:     db["param"],
			TimeStamp: db["ts"].(time.Time),
		}
		return &data
	}
}
