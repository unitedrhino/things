package relationDB

import (
	"github.com/i-Things/things/shared/stores"
)

// 示例
type UdExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

type UdSenseInfo struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`   // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`       // 租户编码
	UserID     int64             `gorm:"column:user_id;type:BIGINT;NOT NULL"`                // 用户id
	Name       string            `gorm:"column:name;type:varchar(100);uniqueIndex;NOT NULL"` // 名称
	Delay      int64             `gorm:"column:delay;type:bigint;NOT NULL"`                  // 延迟
	DelayUnit  string            `gorm:"column:delay_unit;type:varchar(10);NOT NULL"`        // 延迟单位 s:秒 m:分钟 h:小时 d:天
	Desc       string            `gorm:"column:desc;type:varchar(200);NOT NULL"`             // 描述
	HeadImg    string            `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`         // 用户头像

}

type UdSenseDevice struct {
	ID         int64          `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`              // id编号
	SenseID    int64          `gorm:"column:sense_id;type:bigint;NOT NULL"`                          // 信号编号
	ProductID  string         `gorm:"column:product_id;type:char(11);NOT NULL"`                      // 产品编号
	DeviceName string         `gorm:"column:device_name;type:varchar(100);uniqueIndex;NOT NULL"`     // 设备编号
	Params     map[string]any `gorm:"column:params;type:json;serializer:json;NOT NULL;default:'{}'"` // 设备参数
	Delay      int64          `gorm:"column:delay;type:bigint;NOT NULL"`                             // 延迟
	DelayUnit  string         `gorm:"column:delay_unit;type:varchar(10);NOT NULL"`                   // 延迟单位 s:秒 m:分钟 h:小时 d:天
}

type UdAutoInfo struct {
	ID            int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`   // id编号
	TenantCode    stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`       // 租户编码
	UserID        int64             `gorm:"column:user_id;type:BIGINT;NOT NULL"`                // 用户id
	Name          string            `gorm:"column:name;type:varchar(100);uniqueIndex;NOT NULL"` // 名称
	Desc          string            `gorm:"column:desc;type:varchar(200);NOT NULL"`             // 描述
	EffectiveType string            //生效条件关系
	When          []*UdAutoWhen
	EffectiveTime UdAuthEffectiveTime
	Triggers      []*UdAutoTrigger
	Actions       []*UdAutoAction
}
type UdAuthEffectiveTime struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	SenseID    int64             `gorm:"column:sense_id;type:bigint;NOT NULL"`             // 信号编号
	ProductID  string            `gorm:"column:product_id;type:char(11);NOT NULL"`         // 产品编号
	DeviceName string            `gorm:"column:device_name;type:varchar"`
}

type UdAutoTrigger struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	SenseID    int64             `gorm:"column:sense_id;type:bigint;NOT NULL"`             // 信号编号
	AutoID     int64             `gorm:"column:auto_id;type:bigint;NOT NULL"`              // 自动编号
	TriggerID  int64             `gorm:"column:trigger_id;type:bigint;NOT NULL"`           // 触发器编号
}

type UdAutoWhen struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	SenseID    int64             `gorm:"column:sense_id;type:bigint;NOT NULL"`             // 信号编号
}

type UdAutoAction struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	SenseID    int64             `gorm:"column:sense_id;type:bigint;NOT NULL"`             // 信号编号
}
