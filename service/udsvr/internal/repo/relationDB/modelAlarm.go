package relationDB

import (
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"time"
)

// 告警配置信息表
type UdAlarmInfo struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	ProjectID  stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
	Name       string            `gorm:"column:name;type:VARCHAR(100);NOT NULL"`
	Desc       string            `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`
	Level      int64             `gorm:"column:level;type:SMALLINT;default:1"`                   // 告警配置级别（1提醒 2一般 3严重 4紧急 5超紧急）
	Status     int64             `gorm:"column:status;type:SMALLINT;default:1"`                  // 状态 1:启用,2:禁用
	Notifies   []*UdAlarmNotify  `gorm:"column:notifies;type:json;serializer:json;default:'[]'"` // 短信通知模版编码
	stores.Time
}

type UdAlarmNotify struct {
	Type         string   `json:"type"`         //通知类型
	TemplateCode string   `json:"templateCode"` //模版code,不选就是默认的
	UserIDs      []int64  `json:"userIDs"`      //指定用户ID
	Accounts     []string `json:"accounts"`     //账号
}

func (m *UdAlarmInfo) TableName() string {
	return "ud_alarm_info"
}

// 告警配置与场景关联表
type UdAlarmScene struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	ProjectID  stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
	AlarmID    int64             `gorm:"column:alarm_id;uniqueIndex:ai_si;type:BIGINT;NOT NULL"`
	SceneID    int64             `gorm:"column:scene_id;uniqueIndex:ai_si;type:BIGINT;NOT NULL"`
	SceneInfo  *UdSceneInfo      `gorm:"foreignKey:ID;references:SceneID"`
	AlarmInfo  *UdAlarmInfo      `gorm:"foreignKey:ID;references:AlarmID"`
	stores.Time
}

func (m *UdAlarmScene) TableName() string {
	return "ud_alarm_scene"
}

// 告警记录表 一个告警
type UdAlarmRecord struct {
	ID          int64                 `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	TenantCode  stores.TenantCode     `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	ProjectID   stores.ProjectID      `gorm:"column:project_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
	AlarmID     int64                 `gorm:"column:alarm_id;type:BIGINT;NOT NULL"`             //告警记录ID
	AlarmName   string                `gorm:"column:alarm_name;type:VARCHAR(100);NOT NULL"`     //告警名称
	TriggerType scene.TriggerType     `gorm:"column:trigger_type;type:BIGINT;NOT NULL"`         //触发类型(设备触发1,其他2)
	ProductID   string                `gorm:"column:product_id;type:varchar(100);"`             //触发产品id
	DeviceName  string                `gorm:"column:device_name;type:varchar(100);"`            //触发设备名称
	Level       int64                 `gorm:"column:level;type:SMALLINT;NOT NULL"`              //告警配置级别（1提醒 2一般 3严重 4紧急 5超紧急）
	SceneName   string                `gorm:"column:scene_name;type:varchar(100);NOT NULL"`
	SceneID     int64                 `gorm:"column:scene_id;type:BIGINT;NOT NULL"`
	DealStatus  scene.AlarmDealStatus `gorm:"column:deal_status;type:SMALLINT;default:1;NOT NULL"` //告警记录状态（1告警中 2已忽略 3已处理）
	Desc        string                `gorm:"column:desc;type:varchar(100);"`
	WorkOrderID int64                 `gorm:"column:work_order_id;type:BIGINT;NOT NULL"` //工作流ID
	AlarmCount  int64                 `gorm:"column:alarm_count;type:BIGINT;default:1"`  //告警次数
	LastAlarm   time.Time             `gorm:"column:last_alarm;NOT NULL"`
	CreatedTime time.Time             `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *UdAlarmRecord) TableName() string {
	return "ud_alarm_record"
}

//// 告警流水详情表
//type UdAlarmLog struct {
//	ID            int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
//	AlarmRecordID int64     `gorm:"column:alarm_record_id;type:BIGINT;NOT NULL"`
//	Serial        string    `gorm:"column:serial;type:varchar(1024);NOT NULL"`
//	SceneName     string    `gorm:"column:scene_name;type:varchar(100);NOT NULL"`
//	SceneID       int64     `gorm:"column:scene_id;type:BIGINT;NOT NULL"`
//	Desc          string    `gorm:"column:desc;type:varchar(1024);NOT NULL"`
//	CreatedTime   time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"`
//}
//
//func (m *UdAlarmLog) TableName() string {
//	return "ud_alarm_log"
//}

//// 告警处理记录表
//type UdAlarmDealRecord struct {
//	ID            int64     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
//	AlarmRecordID int64     `gorm:"column:alarm_record_id;type:BIGINT;NOT NULL"`
//	Result        string    `gorm:"column:result;type:varchar(1024);NOT NULL"`
//	Type          int64     `gorm:"column:type;type:SMALLINT;NOT NULL"`
//	AlarmTime     time.Time `gorm:"column:alarm_time;default:CURRENT_TIMESTAMP;NOT NULL"`
//	CreatedTime   time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP;NOT NULL"`
//}
//
//func (m *UdAlarmDealRecord) TableName() string {
//	return "ud_alarm_deal_record"
//}
