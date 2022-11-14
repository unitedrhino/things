// Package device 设备发送来的消息解析
package deviceMsg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type (
	PublishMsg struct { //发布消息结构体
		Topic      string
		Payload    []byte
		Timestamp  time.Time
		ProductID  string
		DeviceName string
	}

	CommonMsg struct { //消息内容通用字段
		Method      string `json:"method"`              //操作方法
		ClientToken string `json:"clientToken"`         //方便排查随机数
		Timestamp   int64  `json:"timestamp,omitempty"` //毫秒时间戳
		Code        int64  `json:"code,omitempty"`      //状态码
		Status      string `json:"status,omitempty"`    //返回信息
		Data        any    `json:"data,omitempty"`      //返回具体设备上报的最新数据内容
	}
)

func (p *PublishMsg) String() string {
	msgMap := map[string]any{
		"Topic":       p.Topic,
		"Payload":     string(p.Payload),
		"Timestamp":   p.Timestamp,
		"ProductID":   p.ProductID,
		"DeviceNames": p.DeviceName,
	}
	return utils.Fmt(msgMap)
}

func NewRespCommonMsg(method, clientToken string) *CommonMsg {
	return &CommonMsg{
		Method:      GetRespMethod(method),
		ClientToken: clientToken,
		Timestamp:   time.Now().UnixMilli(),
	}
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
		Topic:      pubInfo.Topic,
		Payload:    pubInfo.Payload,
		Timestamp:  time.UnixMilli(pubInfo.Timestamp),
		ProductID:  pubInfo.ProductID,
		DeviceName: pubInfo.DeviceName,
	}
	return &ele, nil
}

func GenRespTopic(topics any) string {
	var (
		strs []string
		ok   bool
	)
	if strs, ok = topics.([]string); ok {
	} else if str, ok := topics.(string); ok {
		strs = strings.Split(str, "/")
	} else {
		panic("GenRespTopic not support type")
	}
	respTopic := fmt.Sprintf("%s/down/%s/%s/%s",
		strs[0], strs[2], strs[3], strs[4])
	return respTopic
}
