// Package scene 执行动作
package scene

type Action struct {
	Executor string      `json:"executor"` //执行器类型 notify: 通知 delay:延迟  device:设备输出  alarm: 告警
	Delay    ActionDelay `json:"delay""`
	Alarm    ActionAlarm `json:"alarm""`
}

type ActionDelay struct {
	Time int64  `json:"time"` //延迟时间
	Unit string `json:"unit"` //时间单位 seconds:秒  minutes:分钟  hours:小时
}

type ActionAlarm struct {
	Mode string `json:"mode"` //告警模式  trigger: 触发告警  relieve: 解除告警
}
