package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"time"
)

// 设备影子表
type DiDeviceShadow struct {
	ID                int64      `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	ProductID         string     `gorm:"column:product_id;uniqueIndex:pi_dn_di;type:CHAR(11);NOT NULL"`
	DeviceName        string     `gorm:"column:device_name;uniqueIndex:pi_dn_di;type:VARCHAR(100);NOT NULL"`
	DataID            string     `gorm:"column:data_id;uniqueIndex:pi_dn_di;type:VARCHAR(100);NOT NULL"`
	Value             string     `gorm:"column:value;type:VARCHAR(100);default:NULL"`
	UpdatedDeviceTime *time.Time `gorm:"column:updated_device_time;default:NULL"`
	stores.Time
}

func (m *DiDeviceShadow) TableName() string {
	return "di_device_shadow"
}
func (m *DiDeviceShadow) Comment() string {
	return "di_device_shadow"
}
