package tdengine

import (
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"time"
)

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
