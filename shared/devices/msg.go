package devices

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/trace"
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
		Topic      string `json:"topic"` //只用于日志记录
		Timestamp  int64  `json:"timestamp"`
		ProductID  string `json:"productID"`
		DeviceName string `json:"deviceName"`
		Handle     string `json:"handle"` //对应 mqtt topic的第一个 thing ota config 等等
		Type       string `json:"type"`   //操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
		Payload    []byte `json:"payload"`
	}
	// InnerPublish 用于其他服务 发送给ddsvr 转发给设备的
	InnerPublish struct {
		Handle     string `json:"handle"` //对应 mqtt topic的第一个 thing ota config 等等
		Type       string `json:"type"`   // 操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
		Payload    []byte `json:"payload"`
		ProductID  string `json:"productID"`
		DeviceName string `json:"deviceName"`
	}
)

// 发送给设备的数据组包
func PublishToDev(handle string, Type string, payload []byte, productID string, deviceName string) []byte {
	pub := InnerPublish{
		Handle:     handle,
		Type:       Type,
		Payload:    payload,
		ProductID:  productID,
		DeviceName: deviceName,
	}
	data, _ := json.Marshal(pub)
	return data
}

// 收到发送给设备的数据,解包
func GetPublish(data []byte) *InnerPublish {
	pub := InnerPublish{}
	_ = json.Unmarshal(data, &pub)
	return &pub
}

func GenMsgToken(ctx context.Context) string {
	tid := trace.TraceIDFromContext(ctx)
	span := trace.SpanIDFromContext(ctx)
	if tid == "" || span == "" {
		return uuid.NewString()
	}
	return fmt.Sprintf("%s-%s", tid, span)
}
