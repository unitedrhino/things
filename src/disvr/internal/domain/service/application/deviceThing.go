package application

import (
	"github.com/i-Things/things/shared/devices"
)

//属性上报消息体
type PropertyReport struct {
	Device     devices.Core `json:"device"`
	Timestamp  int64        `json:"timestamp,string"` //毫秒时间戳
	Identifier string       `json:"identifier"`       //推送属性的标识符
	Param      any          `json:"param"`            //推送属性的参数
}

type EventReport struct {
	Device     devices.Core   `json:"device"`
	Timestamp  int64          `json:"timestamp,string"` //毫秒时间戳
	Identifier string         `json:"identifier"`       //标识符
	Type       string         `json:"type" `            //事件类型: 信息:info  告警:alert  故障:fault
	Params     map[string]any `json:"params" `          //事件参数
}
