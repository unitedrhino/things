package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"time"
)

// 设备影子表
type DiDeviceShadow struct {
	ID                int64      `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	ProductID         string     `gorm:"column:product_id;type:char(11);NOT NULL"`              // 产品id
	DeviceName        string     `gorm:"column:device_name;type:varchar(100);NOT NULL"`         // 设备名称
	DataID            string     `gorm:"column:data_id;type:varchar(100);NOT NULL"`             // 属性id
	Value             string     `gorm:"column:value;type:varchar(100);default:NULL"`           // 属性值
	UpdatedDeviceTime *time.Time `gorm:"column:updated_device_time;type:datetime;default:NULL"` //更新到设备时间
	stores.Time
}

func (m *DiDeviceShadow) TableName() string {
	return "di_device_shadow"
}
