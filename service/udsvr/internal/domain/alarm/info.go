package alarm

import "gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"

// 告警记录状态
const (
	DealStateNone     = 1 //无告警
	DealStateAlarming = 2 //告警中
	DealStateAlarmed  = 3 //已处理
)

// 告警配置级别
const (
	LevelInfo     = 1 //提醒
	LevelNormal   = 1 //一般
	LevelBad      = 1 //严重
	LevelUrgent   = 1 //紧急
	LevelSpUrgent = 1 //超级紧急
)

type InfoFilter struct {
	Name     string //名字
	SceneID  int64  // 场景ID
	AlarmIDs []int64
}

type Notify struct {
	TriggerType scene.TriggerType     `json:"triggerType,omitempty"` //触发类型(设备触发1,其他2)
	ProductID   string                `json:"productID,omitempty"`   //触发产品id
	DeviceName  string                `json:"deviceName,omitempty"`  //触发设备名称
	SceneID     int64                 `json:"sceneID,omitempty"`     //场景ID
	Mode        scene.ActionAlarmMode `json:"mode,omitempty"`        //告警模式  trigger: 触发告警  relieve: 解除告警
}
