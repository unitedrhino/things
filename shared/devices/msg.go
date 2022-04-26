package devices

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/events"
)

type (
	// DevConn ddsvr 发布设备连接和断连的结构体
	DevConn struct {
		UserName  string `json:"username"`
		Timestamp int64  `json:"ts"` //毫秒时间戳
		Address   string `json:"addr"`
		ClientID  string `json:"clientID"`
		Reason    string `json:"reason"`
		Action    string `json:"action"` //登录 onLogin 登出 onLogout
	}
	// DevPublish ddsvr 发布设备发布信息的结构体
	DevPublish struct {
		Timestamp  int64  `json:"ts"`
		ProductID  string `json:"productID"`
		DeviceName string `json:"deviceName"`
		Topic      string `json:"topic"`
		Payload    []byte `json:"payload"`
	}
	// InnerPublish 用于其他服务发送给ddsvr转发给设备的
	InnerPublish struct {
		Topic   string `json:"topic"`
		Payload []byte `json:"payload"`
	}
)

//发送给设备的数据组包
func PublishToDev(ctx context.Context, topic string, payload []byte) []byte {
	pub := InnerPublish{
		Topic:   topic,
		Payload: payload,
	}
	data, _ := json.Marshal(pub)
	return events.NewEventMsg(ctx, data)
}

//收到发送给设备的数据,解包
func GetPublish(data []byte) (ctx context.Context, topic string, payload []byte) {
	pub := InnerPublish{}
	msg := events.GetEventMsg(data)
	_ = json.Unmarshal(msg.GetData(), &pub)
	return msg.GetCtx(), pub.Topic, pub.Payload
}
