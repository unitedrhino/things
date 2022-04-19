package deviceLog

import "time"

type DeviceLog struct {
	ProductID  string    `json:"productID"`  // 产品id
	DeviceName string    `json:"deviceName"` // 设备名称
	Content    string    `json:"content"`    // 具体信息
	Topic      string    `json:"topic"`      // 主题
	Action     string    `json:"action"`     // 操作类型
	Timestamp  time.Time `json:"timestamp"`  // 操作时间
	RequestID  string    `json:"requestID"`  // 请求ID
	TranceID   string    `json:"tranceID"`   // 服务器端事务id
	ResultType int64     `json:"resultType"` // 请求结果状态,0为成功
}
