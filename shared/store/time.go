package store

import (
	"gorm.io/gorm"
	"time"
)

type Time struct {
	CreatedTime time.Time      `gorm:"column:createdTime;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedTime time.Time      `gorm:"column:updatedTime;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL"`
	DeletedTime gorm.DeletedAt `gorm:"column:deletedTime;type:datetime"`
}
