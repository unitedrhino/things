package relationDB

import "gitee.com/i-Things/share/stores"

type DmUserDeviceCollect struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                              // 租户编码
	ProjectID  stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`                                // 项目ID(雪花ID)
	UserID     int64             `gorm:"column:user_id;type:BIGINT;uniqueIndex:product_id_deviceName;NOT NULL"`           // 问题提出的用户
	ProductID  string            `gorm:"column:product_id;type:char(11);uniqueIndex:product_id_deviceName;NOT NULL"`      // 产品id
	DeviceName string            `gorm:"column:device_name;uniqueIndex:product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:product_id_deviceName"`
}

func (m *DmUserDeviceCollect) TableName() string {
	return "dm_user_device_collect"
}

type DmUserDeviceShare struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                              // 租户编码
	UserID     int64             `gorm:"column:user_id;type:BIGINT;uniqueIndex:product_id_deviceName;NOT NULL"`           // 问题提出的用户
	ProductID  string            `gorm:"column:product_id;type:char(11);uniqueIndex:product_id_deviceName;NOT NULL"`      // 产品id
	DeviceName string            `gorm:"column:device_name;uniqueIndex:product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	AccessPerm []string          `gorm:"column:access_prem;type:json;serializer:json;NOT NULL;default:'[]'"`              //操作权限 hubLog:设备消息记录,ota:ota升级权限,deviceTiming:设备定时
	SchemaPerm []string          `gorm:"column:schema_prem;type:json;serializer:json;NOT NULL;default:'[]'"`              //物模型权限,只需要填写需要授权并授权的物模型id
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:product_id_deviceName"`
}

func (m *DmUserDeviceShare) TableName() string {
	return "dm_user_device_share"
}
