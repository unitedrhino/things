package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"time"
)

// 示例
type RuleExample struct {
	ID int64 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"` // id编号
}

// 规则引擎-场景联动信息表
type RuleSceneInfo struct {
	ID          int64         `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"` // id
	Name        string        `gorm:"column:name;type:varchar(128)"`                        // 场景名称
	TriggerType string        `gorm:"column:triggerType;type:varchar(24);NOT NULL"`         // 触发器类型 device: 设备触发 timer: 定时触发 manual:手动触发
	Trigger     scene.Trigger `gorm:"column:trigger;type:json;serializer:json"`             // 触发器内容-根据触发器类型改变
	When        scene.Terms   `gorm:"column:when;type:json;serializer:json"`                // 触发条件
	Then        scene.Actions `gorm:"column:then;type:json;serializer:json"`                // 满足条件时执行的动作
	Desc        string        `gorm:"column:desc;type:varchar(512)"`                        // 描述
	Status      int64         `gorm:"column:status;type:tinyint(1);default:1"`              // 状态  1:启用,2:禁用
	stores.Time
}

func (m *RuleSceneInfo) TableName() string {
	return "rule_scene_info"
}

// 告警配置与场景关联表
type RuleAlarmScene struct {
	ID      int64 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"` // id编号
	AlarmID int64 `gorm:"column:alarmID;type:bigint(20);NOT NULL"`              // 告警配置ID
	SceneID int64 `gorm:"column:sceneID;type:int(11);NOT NULL"`                 // 场景ID
	stores.Time
}

func (m *RuleAlarmScene) TableName() string {
	return "rule_alarm_scene"
}

// 告警配置信息表
type RuleAlarmInfo struct {
	ID     int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"` // 编号
	Name   string `gorm:"column:name;type:varchar(100);NOT NULL"`               // 告警配置名称
	Desc   string `gorm:"column:desc;type:varchar(100);NOT NULL"`               // 告警配置说明
	Level  int64  `gorm:"column:level;type:tinyint(1);NOT NULL"`                // 告警配置级别（1提醒 2一般 3严重 4紧急 5超紧急）
	Status int64  `gorm:"column:status;type:tinyint(1);default:1"`              // 状态  1:启用,2:禁用
	stores.Time
}

func (m *RuleAlarmInfo) TableName() string {
	return "rule_alarm_info"
}

// 告警记录表
type RuleAlarmRecord struct {
	ID          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`                // 编号
	AlarmID     int64     `gorm:"column:alarmID;type:bigint(20);NOT NULL"`                             // 告警记录ID
	TriggerType int64     `gorm:"column:triggerType;type:int(11);NOT NULL"`                            // 触发类型(设备触发1,其他2)
	ProductID   string    `gorm:"column:productID;type:char(11);NOT NULL"`                             // 触发产品id
	DeviceName  string    `gorm:"column:deviceName;type:varchar(100);NOT NULL"`                        // 触发设备名称
	Level       int64     `gorm:"column:level;type:tinyint(1);NOT NULL"`                               // 告警配置级别（1提醒 2一般 3严重 4紧急 5超紧急）
	SceneName   string    `gorm:"column:sceneName;type:varchar(100);NOT NULL"`                         // 场景名称
	SceneID     int64     `gorm:"column:sceneID;type:int(11);NOT NULL"`                                // 场景ID
	DealState   int64     `gorm:"column:dealState;type:tinyint(1);default:1;NOT NULL"`                 // 告警记录状态（1无告警 2告警中 3已处理）
	LastAlarm   time.Time `gorm:"column:lastAlarm;type:datetime;NOT NULL"`                             // 最新告警时间
	CreatedTime time.Time `gorm:"column:createdTime;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL"` // 告警时间
}

func (m *RuleAlarmRecord) TableName() string {
	return "rule_alarm_record"
}

// 告警流水详情表
type RuleAlarmLog struct {
	ID            int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`                // 编号
	AlarmRecordID int64     `gorm:"column:alarmRecordID;type:bigint(20);NOT NULL"`                       // 告警记录ID
	Serial        string    `gorm:"column:serial;type:varchar(1024);NOT NULL"`                           // 告警流水
	SceneName     string    `gorm:"column:sceneName;type:varchar(100);NOT NULL"`                         // 场景名称
	SceneID       int64     `gorm:"column:sceneID;type:int(11);NOT NULL"`                                // 场景ID
	Desc          string    `gorm:"column:desc;type:varchar(1024);NOT NULL"`                             // 告警说明
	CreatedTime   time.Time `gorm:"column:createdTime;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL"` // 告警时间
}

func (m *RuleAlarmLog) TableName() string {
	return "rule_alarm_log"
}

// 告警处理记录表
type RuleAlarmDealRecord struct {
	ID            int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`                // 编号
	AlarmRecordID int64     `gorm:"column:alarmRecordID;type:bigint(20);NOT NULL"`                       // 告警记录ID
	Result        string    `gorm:"column:result;type:varchar(1024);NOT NULL"`                           // 告警处理结果
	Type          int64     `gorm:"column:type;type:tinyint(1);NOT NULL"`                                // 告警处理类型（1人工 2系统）
	AlarmTime     time.Time `gorm:"column:alarmTime;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL"`   // 告警时间
	CreatedTime   time.Time `gorm:"column:createdTime;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL"` // 告警处理时间
}

func (m *RuleAlarmDealRecord) TableName() string {
	return "rule_alarm_deal_record"
}
