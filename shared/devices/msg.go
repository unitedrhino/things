package devices

import (
	"encoding/json"
)

type (
	// DevConn ddsvr 发布设备 连接和断连 的结构体
	DevConn struct {
		UserName  string `json:"username"`
		Timestamp int64  `json:"timestamp"` //毫秒时间戳
		Address   string `json:"addr"`
		ClientID  string `json:"clientID"`
		Reason    string `json:"reason"`
		Action    string `json:"action"` //登录 onLogin 登出 onLogout
	}
	// DevPublish ddsvr 发布设备 发布信息 的结构体
	DevPublish struct {
		Timestamp  int64  `json:"timestamp"`
		ProductID  string `json:"productID"`
		DeviceName string `json:"deviceName"`
		Topic      string `json:"topic"`
		Payload    []byte `json:"payload"`
	}
	// InnerPublish 用于其他服务 发送给ddsvr 转发给设备的
	InnerPublish struct {
		Topic   string `json:"topic"`
		Payload []byte `json:"payload"`
	}
)

//发送给设备的数据组包
func PublishToDev(topic string, payload []byte) []byte {
	pub := InnerPublish{
		Topic:   topic,
		Payload: payload,
	}
	data, _ := json.Marshal(pub)
	return data
}

//收到发送给设备的数据,解包
func GetPublish(data []byte) (topic string, payload []byte) {
	pub := InnerPublish{}
	_ = json.Unmarshal(data, &pub)
	return pub.Topic, pub.Payload
}
