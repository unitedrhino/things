package deviceTimer

import (
	"gitee.com/i-Things/share/devices"
)

type Info struct {
	ID          int64        `json:"id,omitempty"`   //场景id
	Name        string       `json:"name,omitempty"` //名称
	Device      devices.Core `json:"device,omitempty"`
	CreatedTime int64        `json:"createdTime,omitempty"` //创建时间 秒级时间戳 只读
	TriggerType Trigger      `json:"triggerType,omitempty"` //触发类型 timer: 定时触发 delay: 延迟触发(延迟触发同时只能存在一个) timeRange:时间段触发
	Repeat      int64        `json:"repeat,omitempty"`      //重复 二进制周一到周日 11111111 这个参数只有定时触发才有
	ExecAt      int64        `json:"execAt,omitempty"`      //执行时间 从0点加起来的秒数 如 1点就是 1*60*60
	ExecEndAt   int64        `json:"execAtEnd,omitempty"`   //结束执行时间 从0点加起来的秒数 如 1点就是 1*60*60
	Exec        Action       `json:"exec"`                  //执行动作
	ExecEnd     *Action      `json:"execEnd"`               //结束执行动作
	LastRunTime int64        `json:"lastRunTime,omitempty"` //最后一次执行时间 秒级时间戳 只读
	Status      int64        `json:"status,omitempty"`      // 状态（1启用 2禁用）
}

type Action struct {
	ActionType string `json:"actionType,omitempty"` //云端向设备发起属性控制: propertyControl  应用调用设备行为:action
	DataID     string `json:"dataID,omitempty"`     //属性的id及事件的id
	Value      string `json:"value,omitempty"`      //传的值
}

type Trigger = string

const (
	TriggerTimer     Trigger = "timer"
	TriggerDelay     Trigger = "delay"
	TriggerTimeRange Trigger = "timeRange"
)
