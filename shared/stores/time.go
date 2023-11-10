package stores

import (
	"database/sql"
	"time"
)

type Time struct {
	CreatedTime time.Time    `gorm:"column:created_time;index;sort:desc;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedTime time.Time    `gorm:"column:updated_time;autoUpdateTime;default:CURRENT_TIMESTAMP;NOT NULL"`
	DeletedTime sql.NullTime `gorm:"column:deleted_time"`
}

type RecordTime struct {
	CreatedTime time.Time    `gorm:"column:created_time;index:,sort:desc;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedTime time.Time    `gorm:"column:updated_time;autoUpdateTime;default:CURRENT_TIMESTAMP;NOT NULL"`
	RecordDate  sql.NullTime `gorm:"column:record_date"`
}
