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
	PublishMsg = devices.DevPublish

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
		c.Msg = e.GetI18nMsg("")
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
	pubInfo := PublishMsg{}
	err := json.Unmarshal(data, &pubInfo)
	if err != nil {
		logx.WithContext(ctx).Error("GetDevPublish", string(data), err)
		return nil, err
	}
	ele := utils.Copy[PublishMsg](pubInfo)
	return ele, nil
}
