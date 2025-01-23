// Package device 设备发送来的消息解析
package deviceMsg

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ReqType = string

const (
	ReqMsg  ReqType = "req"
	RespMsg ReqType = "resp"
)

type (
	PublishMsg struct { //发布消息结构体
		Topic        string            `json:"topic"`  //只用于日志记录
		Handle       devices.MsgHandle `json:"handle"` //对应 mqtt topic的第一个 thing ota config 等等
		Type         string            `json:"type"`   //操作类型 从topic中提取 物模型下就是   property属性 event事件 action行为
		Payload      []byte            `json:"payload"`
		Timestamp    int64             `json:"timestamp"` //毫秒时间戳
		ProductID    string            `json:"productID"`
		DeviceName   string            `json:"deviceName"`
		Explain      string            `json:"explain"`      //内部使用的拓展字段
		ProtocolCode string            `json:"protocolCode"` //如果有该字段则回复的时候也会带上该字段
	}

	CommonMsg struct { //消息内容通用字段
		Method    string     `json:"method,omitempty"`    //操作方法
		MsgToken  string     `json:"msgToken,omitempty"`  //方便排查随机数
		Timestamp int64      `json:"timestamp,omitempty"` //毫秒时间戳
		Code      int64      `json:"code,omitempty"`      //状态码
		Msg       string     `json:"msg,omitempty"`       //返回信息
		Data      any        `json:"data,omitempty"`      //返回具体设备上报的最新数据内容
		Sys       *SysConfig `json:"sys,omitempty"`       //系统配置
	}
	SysConfig struct {
		NoAsk  bool `json:"noAsk"`  //云平台是否回复消息
		RetMsg bool `json:"retMsg"` //是否返回错误信息
	}
)

func (c *CommonMsg) NeedRetMsg() bool {
	if c.Sys != nil {
		return c.Sys.RetMsg
	}
	return false
}

func (p *PublishMsg) String() string {
	msgMap := map[string]any{
		"Handle":       p.Handle,
		"Type":         p.Type,
		"Payload":      string(p.Payload),
		"Timestamp":    p.Timestamp,
		"ProductID":    p.ProductID,
		"DeviceName":   p.DeviceName,
		"protocolCode": p.ProtocolCode,
	}
	return utils.Fmt(msgMap)
}

func (p *PublishMsg) GetPayload() string {
	if p == nil || len(p.Payload) == 0 {
		return ""
	}
	return string(p.Payload)
}

// 如果MsgToken为空,会使用uuid生成一个
func NewRespCommonMsg(ctx context.Context, method, MsgToken string) *CommonMsg {
	return &CommonMsg{
		Method:   GetRespMethod(method),
		MsgToken: MsgToken,
		//Timestamp: time.Now().UnixMilli(),
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
func (c *CommonMsg) AddStatus(err error, needRet bool) *CommonMsg {
	if err == nil {
		err = errors.OK
	}
	e := errors.Fmt(err)
	c.Code = e.GetCode()
	if needRet {
		c.Msg = e.GetI18nMsg("en")
	}
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
	ele := utils.Copy[PublishMsg](pubInfo)
	return ele, nil
}
