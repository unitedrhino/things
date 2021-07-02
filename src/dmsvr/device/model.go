package device

import (
	"time"
)

type DeviceData struct {
	Property  map[string]interface{}   `json:"property"` //属性
	Event     map[string]interface{}   `json:"event"` //事件
	TimeStamp time.Time `json:"timeStamp"`//时间戳
}
