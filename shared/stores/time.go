package stores

import (
	"database/sql"
	"time"
)

type Time struct {
	CreatedTime time.Time    `gorm:"column:created_time;type:timestamp without time zone;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedTime time.Time    `gorm:"column:updated_time;type:timestamp without time zone;default:CURRENT_TIMESTAMP;NOT NULL"`
	DeletedTime sql.NullTime `gorm:"column:deleted_time;type:timestamp without time zone"`
}
