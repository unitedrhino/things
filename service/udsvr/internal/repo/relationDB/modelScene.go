package relationDB

import (
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"time"
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
	LastRunTime    time.Time         `gorm:"column:last_run_time;index;default:CURRENT_TIMESTAMP;NOT NULL"`
	Status         int64             `gorm:"column:status;type:BIGINT;default:1"` //状态
	UdSceneTrigger `gorm:"embedded;embeddedPrefix:trigger_"`
	UdSceneWhen    `gorm:"embedded;embeddedPrefix:when_"`
	UdSceneThen    `gorm:"embedded;embeddedPrefix:then_"`
	stores.SoftTime
}

func (m *UdSceneInfo) TableName() string {
	return "ud_scene_info"
}

type UdSceneTrigger struct {
	Type    string                  `gorm:"column:type;type:VARCHAR(25);NOT NULL"` //触发类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Devices []*UdSceneTriggerDevice `gorm:"foreignKey:SceneID;references:ID"`
	Timers  []*UdSceneTriggerTimer  `gorm:"foreignKey:SceneID;references:ID"`
}

type UdSceneTriggerTimer struct {
	ID          int64        `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`       // id编号
	SceneID     int64        `gorm:"column:scene_id;index;type:bigint"`                      // 场景id编号
	ExecAt      int64        `gorm:"column:exec_at;index;type:bigint;NOT NULL"`              //执行时间 从0点加起来的秒数 如 1点就是 1*60*60
	ExecRepeat  int64        `gorm:"column:exec_repeat;index;type:bigint;default:0b1111111"` //重复 二进制周日到周六 11111111 这个参数只有定时触发才有
	LastRunTime time.Time    `gorm:"column:last_run_time;index;default:CURRENT_TIMESTAMP;NOT NULL"`
	SceneInfo   *UdSceneInfo `gorm:"foreignKey:ID;references:SceneID"`
	Status      int64        `gorm:"column:status;type:BIGINT;default:1"` //状态 同步场景联动的status
	stores.Time
}

func (m *UdSceneTriggerTimer) TableName() string {
	return "ud_scene_trigger_timer"
}

type UdSceneTriggerDevice struct {
	ID              int64                  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`  // id编号
	SceneID         int64                  `gorm:"column:scene_id;index;type:bigint"`                 // 场景id编号
	ProductID       string                 `gorm:"column:product_id;index;type:VARCHAR(25);NOT NULL"` //产品id
	Selector        string                 `gorm:"column:selector;type:VARCHAR(25);NOT NULL"`         //设备选择方式  all: 全部 fixed:指定的设备
	SelectorValues  []string               `gorm:"column:selector_values;type:json;serializer:json"`  //选择的列表  选择的列表, fixed类型是设备名列表
	Operator        string                 `gorm:"column:operator;type:VARCHAR(25);NOT NULL"`         //触发类型  connected:上线 disConnected:下线 reportProperty:属性上报 reportEvent: 事件上报
	OperationSchema *scene.OperationSchema `gorm:"column:operation_schema;type:json;serializer:json"` //物模型类型的具体操作 reportProperty:属性上报 reportEvent: 事件上报
	SceneInfo       *UdSceneInfo           `gorm:"foreignKey:ID;references:SceneID"`
	stores.Time
}

func (m *UdSceneTriggerDevice) TableName() string {
	return "ud_scene_trigger_device"
}

type UdSceneWhen struct {
	ValidRanges   scene.WhenRanges `gorm:"column:validRanges;type:json;serializer:json"`
	InvalidRanges scene.WhenRanges `gorm:"column:invalidRanges;type:json;serializer:json"`
	Conditions    scene.Conditions `gorm:"column:conditions;type:json;serializer:json"`
}

type UdSceneThen struct {
	Actions scene.Actions `gorm:"column:actions;type:json;serializer:json"`
}
