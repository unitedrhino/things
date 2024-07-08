package relationDB

import (
	"database/sql"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
)

type DmUserDeviceCollect struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                              // 租户编码
	ProjectID  stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`                                // 项目ID(雪花ID)
	UserID     int64             `gorm:"column:user_id;type:BIGINT;uniqueIndex:product_id_deviceName;NOT NULL"`           // 问题提出的用户
	ProductID  string            `gorm:"column:product_id;type:varchar(100);uniqueIndex:product_id_deviceName;NOT NULL"`  // 产品id
	DeviceName string            `gorm:"column:device_name;uniqueIndex:product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:product_id_deviceName"`
}

func (m *DmUserDeviceCollect) TableName() string {
	return "dm_user_device_collect"
}

const (
	ShareAuthTypeAll = 1 //授予全部权限
	ShareAuthType    = 1 //授予全部权限
)

type DmUserDeviceShare struct {
	ID                int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode        stores.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`                           // 租户编码
	SharedUserID      int64             `gorm:"column:shared_user_id;type:BIGINT;uniqueIndex:product_id_deviceName;NOT NULL"` // 分享对象的用户ID
	SharedUserAccount string            `gorm:"column:shared_user_account;type:VARCHAR(100);"`                                // 分享对象的用户账号

	ProjectID  int64                 `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`                                // 分享的设备所属的项目
	ProductID  string                `gorm:"column:product_id;type:varchar(100);uniqueIndex:product_id_deviceName;NOT NULL"`  // 产品id
	DeviceName string                `gorm:"column:device_name;uniqueIndex:product_id_deviceName;type:varchar(100);NOT NULL"` // 设备名称
	AuthType   def.AuthType          `gorm:"column:auth_type;type:varchar(100);default:1"`                                    // 权限类型
	AccessPerm map[string]*SharePerm `gorm:"column:access_prem;type:json;serializer:json;NOT NULL;default:'{}'"`              //操作权限 hubLog:设备消息记录,ota:ota升级权限,deviceTiming:设备定时
	SchemaPerm map[string]*SharePerm `gorm:"column:schema_prem;type:json;serializer:json;NOT NULL;default:'{}'"`              //物模型权限,只需要填写需要授权并授权的物模型id
	ExpTime    sql.NullTime          `gorm:"column:exp_time"`                                                                 //过期时间,为0不限制

	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:product_id_deviceName"`
}
type SharePerm struct {
	Perm int64 `json:"perm"` //1:r(只读) 3(默认):rw(可读可写)
}

func (m *DmUserDeviceShare) TableName() string {
	return "dm_user_device_share"
}
