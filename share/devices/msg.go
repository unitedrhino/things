package devices

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"unicode"
	"unicode/utf8"

	"gitee.com/unitedrhino/share/utils"
)

type Action = string

const (
	ActionConnected    Action = "connected"
	ActionDisconnected Action = "disconnected"
)

type (
	// DevConn ddsvr 发布设备 连接和断连 的结构体
	DevConn struct {
		UserName  string `json:"username"`
		Timestamp int64  `json:"timestamp"` //毫秒时间戳
		Address   string `json:"addr"`
		ClientID  string `json:"clientID"`
		/*
			https://www.emqx.com/en/blog/emqx-mqtt-broker-connection-troubleshooting
				normal：客户端主动断开连接；
				kicked：被服务器踢出，通过 REST API；
				keepalive_timeout：保持活动超时；
				not_authorized：认证失败，或者acl_nomatch=disconnect时，没有权限的Pub/Sub会主动断开客户端连接；
				tcp_closed：对端关闭了网络连接；
				discard：因为相同ClientID的客户端上线了，并且设置了clean_start=true；
				takeovered：因为相同ClientID的客户端上线，并且设置了clean_start=false；
				internal_error：格式错误的消息或其他未知错误。
		*/
		Reason     string `json:"reason"`
		Action     Action `json:"action"` //登录 connected 登出 disconnected
		ProductID  string `json:"productID"`
		DeviceName string `json:"deviceName"`
	}
	// DevPublish ddsvr 发布设备 发布信息 的结构体
	DevPublish struct { //发布消息结构体
		Topic        string    `json:"topic"`  //只用于日志记录
		Handle       MsgHandle `json:"handle"` //对应 mqtt topic的第一个 thing ota config 等等
		Type         string    `json:"type"`   //操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
		Payload      []byte    `json:"payload"`
		Timestamp    int64     `json:"timestamp"` //毫秒时间戳
		ProductID    string    `json:"productID"`
		DeviceName   string    `json:"deviceName"`
		Explain      string    `json:"explain"`      //内部使用的拓展字段
		ProtocolCode string    `json:"protocolCode"` //如果有该字段则回复的时候也会带上该字段
	}
	// InnerPublish 用于其他服务 发送给ddsvr 转发给设备的
	InnerPublish = DevPublish
)

// IsLikelyText 判断字节切片是否更可能是文本
func IsLikelyText(b []byte) bool {
	validUTF8 := utf8.Valid(b)
	if !validUTF8 {
		return false
	}
	total := len(b)
	nonPrintable := 0
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if !unicode.IsPrint(r) {
			nonPrintable++
		}
		b = b[size:]
	}
	// 如果不可打印字符占比超过 20%，则认为是二进制数据
	return float64(nonPrintable)/float64(total) < 0.2
}

func printBytes(data []byte) string {
	// 检查是否为有效的 UTF-8 字符串
	if IsLikelyText(data) {
		// 如果是字符串，直接打印字符串
		return string(data)
	} else {
		// 如果是二进制数据，打印十六进制格式
		return "0x" + hex.EncodeToString(data)
	}
}

func (p *DevPublish) String() string {
	msgMap := map[string]any{
		"Handle":       p.Handle,
		"Type":         p.Type,
		"Payload":      printBytes(p.Payload),
		"Timestamp":    p.Timestamp,
		"ProductID":    p.ProductID,
		"DeviceName":   p.DeviceName,
		"protocolCode": p.ProtocolCode,
	}
	return utils.Fmt(msgMap)
}

func (p *DevPublish) GetPayload() string {
	if p == nil || len(p.Payload) == 0 {
		return ""
	}
	return string(p.Payload)
}

// 发送给设备的数据组包
func PublishToDev(handle string, Type string, payload []byte, protocolCode string, productID string, deviceName string) []byte {
	pub := InnerPublish{
		Handle:       handle,
		Type:         Type,
		Payload:      payload,
		ProtocolCode: protocolCode,
		ProductID:    productID,
		DeviceName:   deviceName,
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

var randID atomic.Uint32

func GenMsgToken(ctx context.Context, nodeID int64) string {
	var token = uint32(nodeID) & 0xff
	token += randID.Add(1) << 8 & 0xfff00
	return fmt.Sprintf("%x", token)
}
