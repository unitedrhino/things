package device

import (
	"time"
)

type DeviceData struct {
	Property map[string]interface{} `json:"property"` //属性
	Event    struct {
		ID     string                 `json:"id"`     //事件id
		Type   string                 `json:"type"`   //事件类型: 信息:info  告警alert  故障:fault
		Params map[string]interface{} `json:"params"` //事件参数
	} `json:"event"` //事件
	TimeStamp time.Time `json:"timeStamp"` //时间戳
}
