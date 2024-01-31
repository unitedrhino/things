package relationDB

import (
	"gitee.com/i-Things/core/shared/stores"
	"github.com/i-Things/things/service/rulesvr/internal/domain/scene"
	"time"
)

// 示例
type RuleExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 规则引擎-场景联动信息表
type RuleSceneInfo struct {
	ID          int64         `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	Name        string        `gorm:"column:name;uniqueIndex;type:varchar(128)"`
	TriggerType string        `gorm:"column:trigger_type;index;type:varchar(24);NOT NULL"`
	Trigger     scene.Trigger `gorm:"column:trigger;type:json;serializer:json"`
	When        scene.Terms   `gorm:"column:when;type:json;serializer:json"`
	Then        scene.Actions `gorm:"column:then;type:json;serializer:json"`
	Desc        string        `gorm:"column:desc;type:varchar(512)"`
	Status      int64         `gorm:"column:status;index;type:BIGINT;default:1"`
	stores.Time
}

func (m *RuleSceneInfo) TableName() string {
	return "rule_scene_info"
}

// 告警配置与场景关联表
type RuleAlarmScene struct {
	ID      int64 `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	AlarmID int64 `gorm:"column:alarm_id;uniqueIndex:ai_si;type:BIGINT;NOT NULL"`
	SceneID int64 `gorm:"column:scene_id;uniqueIndex:ai_si;type:BIGINT;NOT NULL"`
	stores.Time
}

func (m *RuleAlarmScene) TableName() string {
	return "rule_alarm_scene"
}

// 告警配置信息表
type RuleAlarmInfo struct {
	ID     int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	Name   string `gorm:"column:name;uniqueIndex;type:VARCHAR(100);NOT NULL"`
	Desc   string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`
	Level  int64  `gorm:"column:level;type:SMALLINT;NOT NULL"`   // 告警配置级别（1提醒 2一般 3严重 4紧急 5超紧急）
	Status int64  `gorm:"column:status;type:SMALLINT;default:1"` // 状态 1:启用,2:禁用
	stores.Time
}

func (m *RuleAlarmInfo) TableName() string {
	return "rule_alarm_info"
}

// 告警记录表
type RuleAlarmRecord struct {
	ID          int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	AlarmID     int64     `gorm:"column:alarm_id;type:BIGINT;NOT NULL"`
	TriggerType int64     `gorm:"column:trigger_type;uniqueIndex:tt_pi_dn;type:BIGINT;NOT NULL"`
	ProductID   string    `gorm:"column:product_id;uniqueIndex:tt_pi_dn;type:char(11);NOT NULL"`
	DeviceName  string    `gorm:"column:device_name;uniqueIndex:tt_pi_dn;type:varchar(100);NOT NULL"`
	Level       int64     `gorm:"column:level;type:SMALLINT;NOT NULL"`
	SceneName   string    `gorm:"column:scene_name;type:varchar(100);NOT NULL"`
	SceneID     int64     `gorm:"column:scene_id;type:BIGINT;NOT NULL"`
	DealState   int64     `gorm:"column:deal_state;type:SMALLINT;default:1;NOT NULL"`
	LastAlarm   time.Time `gorm:"column:last_alarm;NOT NULL"`
	CreatedTime time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *RuleAlarmRecord) TableName() string {
	return "rule_alarm_record"
}

// 告警流水详情表
type RuleAlarmLog struct {
	ID            int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	AlarmRecordID int64     `gorm:"column:alarm_record_id;type:BIGINT;NOT NULL"`
	Serial        string    `gorm:"column:serial;type:varchar(1024);NOT NULL"`
	SceneName     string    `gorm:"column:scene_name;type:varchar(100);NOT NULL"`
	SceneID       int64     `gorm:"column:scene_id;type:BIGINT;NOT NULL"`
	Desc          string    `gorm:"column:desc;type:varchar(1024);NOT NULL"`
	CreatedTime   time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *RuleAlarmLog) TableName() string {
	return "rule_alarm_log"
}

// 告警处理记录表
type RuleAlarmDealRecord struct {
	ID            int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	AlarmRecordID int64     `gorm:"column:alarm_record_id;type:BIGINT;NOT NULL"`
	Result        string    `gorm:"column:result;type:varchar(1024);NOT NULL"`
	Type          int64     `gorm:"column:type;type:SMALLINT;NOT NULL"`
	AlarmTime     time.Time `gorm:"column:alarm_time;default:CURRENT_TIMESTAMP;NOT NULL"`
	CreatedTime   time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *RuleAlarmDealRecord) TableName() string {
	return "rule_alarm_deal_record"
}
