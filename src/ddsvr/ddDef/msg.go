package ddDef

type (
	DevConn struct {
		UserName  string `json:"username"`
		Timestamp int64  `json:"ts"`
		Address   string `json:"addr"`
		ClientID  string `json:"clientID"`
		Reason    string `json:"reason"`
		Action    string `json:"action"` //登录 onLogin 登出 onLogout
	}
	DevPublish struct {
		Timestamp int64  `json:"ts"`
		Topic     string `json:"topic"`
		Payload   []byte `json:"payload"`
	}
	InnerPublish struct {
		Topic   string `json:"topic"`
		Payload string `json:"payload"`
	}
)

const (
	ActionLogin  = "onLogin"
	ActionLogout = "onLogout"
)

//topic 定义
const (
	// TopicDevPublish dd模块收到设备的发布消息后向内部推送以下topic
	TopicDevPublish = "dd.thing.device.clients.publish"
	// TopicDevConnected dd模块收到设备的登录消息后向内部推送以下topic
	TopicDevConnected = "dd.thing.device.clients.connected"
	// TopicDevDisconnected dd模块收到设备的登出消息后向内部推送以下topic
	TopicDevDisconnected = "dd.thing.device.clients.disconnected"
	// TopicInnerPublish dd模块订阅以下topic,收到内部的发布消息后向设备推送
	TopicInnerPublish = "dd.thing.inner.publish"
)
