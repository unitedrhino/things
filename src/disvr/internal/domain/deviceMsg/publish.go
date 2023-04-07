// Package device 设备发送来的消息解析
package deviceMsg

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

const (
	ReqMsg  = "req"
	RespMsg = "resp"
)

type (
	PublishMsg struct { //发布消息结构体
		Topic      string //只用于日志记录
		Handle     string //对应 mqtt topic的第一个 thing ota config 等等
		Type       string //操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
		Payload    []byte
		Timestamp  int64 //毫秒时间戳
		ProductID  string
		DeviceName string
	}

	CommonMsg struct { //消息内容通用字段
		Method      string     `json:"method"`              //操作方法
		ClientToken string     `json:"clientToken"`         //方便排查随机数
		Timestamp   int64      `json:"timestamp,omitempty"` //毫秒时间戳
		Code        int64      `json:"code,omitempty"`      //状态码
		Status      string     `json:"status,omitempty"`    //返回信息
		Data        any        `json:"data,omitempty"`      //返回具体设备上报的最新数据内容
		Sys         *SysConfig `json:"sys,omitempty"`       //系统配置
	}
	SysConfig struct {
		NoAsk bool `json:"noAsk"` //云平台是否回复消息
	}
)

func (p *PublishMsg) String() string {
	msgMap := map[string]any{
		"Handle":      p.Handle,
		"Type":        p.Type,
		"Payload":     string(p.Payload),
		"Timestamp":   p.Timestamp,
		"ProductID":   p.ProductID,
		"DeviceNames": p.DeviceName,
	}
	return utils.Fmt(msgMap)
}

// 如果clientToken为空,会使用uuid生成一个
func NewRespCommonMsg(method, clientToken string) *CommonMsg {
	if clientToken == "" {
		clientToken, _ = uuid.GenerateUUID()
	}
	return &CommonMsg{
		Method:      GetRespMethod(method),
		ClientToken: clientToken,
		Timestamp:   time.Now().UnixMilli(),
	}
}
func (c *CommonMsg) NoAsk() bool {
	if c.Sys == nil {
		return false
	}
	return c.Sys.NoAsk
}

func (c *CommonMsg) GetTimeStamp() time.Time {
	if c.Timestamp != 0 {
		return time.UnixMilli(c.Timestamp)
	}
	return time.Now()
}
func (c *CommonMsg) AddStatus(err error) *CommonMsg {
	e := errors.Fmt(err)
	c.Code = e.Code
	c.Status = e.GetDetailMsg()
	return c
}
func (c *CommonMsg) Bytes() []byte {
	str, _ := json.Marshal(c)
	return str
}

func (c *CommonMsg) String() string {
	return string(c.Bytes())
}

func GetDevPublish(ctx context.Context, data []byte) (*PublishMsg, error) {
	pubInfo := devices.DevPublish{}
	err := json.Unmarshal(data, &pubInfo)
	if err != nil {
		logx.WithContext(ctx).Error("GetDevPublish", string(data), err)
		return nil, err
	}
	ele := PublishMsg{
		Handle:     pubInfo.Handle,
		Topic:      pubInfo.Topic,
		Type:       pubInfo.Type,
		Payload:    pubInfo.Payload,
		Timestamp:  pubInfo.Timestamp,
		ProductID:  pubInfo.ProductID,
		DeviceName: pubInfo.DeviceName,
	}
	return &ele, nil
}
