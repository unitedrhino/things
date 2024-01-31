package relationDB

import (
	"gitee.com/i-Things/core/shared/stores"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
)

type UdSceneInfo struct {
	ID             int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`     // id编号
	TenantCode     stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`         // 租户编码
	ProjectID      stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`     // 项目ID(雪花ID)
	Tag            string            `gorm:"column:tag;type:VARCHAR(128);NOT NULL;default:normal"` //标签 admin: 管理员 normal: 普通
	AreaIDs        []int64           `gorm:"column:area_ids;type:json;serializer:json"`            // 涉及到的区域列表(需要鉴权)
	HeadImg        string            `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`           // 头像
	Name           string            `gorm:"column:name;type:varchar(100);NOT NULL"`               // 名称
	Desc           string            `gorm:"column:desc;type:varchar(200);NOT NULL"`               // 描述
	Status         int64             `gorm:"column:status;type:BIGINT;default:1"`                  //状态
	UdSceneTrigger `gorm:"embedded;embeddedPrefix:trigger_"`
	UdSceneWhen    `gorm:"embedded;embeddedPrefix:when_"`
	UdSceneThen    `gorm:"embedded;embeddedPrefix:then_"`
	stores.SoftTime
}

func (m *UdSceneInfo) TableName() string {
	return "ud_scene_info"
}

type UdSceneTrigger struct {
	Type    string               `gorm:"column:type;type:VARCHAR(25);NOT NULL"` //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Devices scene.TriggerDevices `gorm:"column:devices;type:json;serializer:json"`
	Timers  scene.Timers         `gorm:"column:timers;type:json;serializer:json"`
}

type UdSceneWhen struct {
	ValidRanges   scene.WhenRanges `gorm:"column:validRanges;type:json;serializer:json"`
	InvalidRanges scene.WhenRanges `gorm:"column:invalidRanges;type:json;serializer:json"`
	Conditions    scene.Conditions `gorm:"column:conditions;type:json;serializer:json"`
}

type UdSceneThen struct {
	Actions scene.Actions `gorm:"column:actions;type:json;serializer:json"`
}
