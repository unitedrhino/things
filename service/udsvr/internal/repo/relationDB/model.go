package relationDB

import (
	"gitee.com/i-Things/share/stores"
	"time"
)

// 示例
type UdExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

type UdDeviceTimerInfo struct {
	ID          int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	TenantCode  stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
	ProjectID   stores.ProjectID  `gorm:"column:project_id;type:bigint;default:0;NOT NULL"` // 项目ID(雪花ID)
	ProductID   string            `gorm:"column:product_id;type:varchar(100);NOT NULL"`     // 产品id
	DeviceName  string            `gorm:"column:device_name;type:varchar(100);NOT NULL"`    // 设备名称
	DeviceAlias string            `gorm:"column:device_alias;type:varchar(100);NOT NULL"`   // 设备名称
	TriggerType string            `gorm:"column:trigger_type;type:varchar(25);NOT NULL"`    //触发类型 timer: 定时触发 delay: 延迟触发(延迟触发同时只能存在一个)
	ExecAt      int64             `gorm:"column:exec_at;type:bigint;NOT NULL"`              //执行时间 从0点加起来的秒数 如 1点就是 1*60*60
	ExecRepeat  int64             `gorm:"column:exec_repeat;type:bigint;default:0b1111111"` //重复 二进制周日到周六 11111111 这个参数只有定时触发才有
	ActionType  string            `gorm:"column:action_type;type:varchar(25);NOT NULL"`     //云端向设备发起属性控制: propertyControl  应用调用设备行为:action
	DataName    string            `gorm:"column:data_name;type:VARCHAR(500);NOT NULL"`      //对应的物模型定义,只读
	DataID      string            `gorm:"column:data_id;type:varchar(100);NOT NULL"`        //属性的id及行为的id
	Value       string            `gorm:"column:value;type:varchar(1024);default:NULL"`     //传的值
	Name        string            `gorm:"column:name;type:varchar(100);default:''"`         // 名称
	LastRunTime time.Time         `gorm:"column:last_run_time;index;default:CURRENT_TIMESTAMP;NOT NULL"`
	Status      int64             `gorm:"column:status;type:BIGINT;default:1"` //状态
	stores.SoftTime
}

func (m *UdDeviceTimerInfo) TableName() string {
	return "ud_device_timer_info"
}
