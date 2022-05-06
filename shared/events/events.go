package events

import (
	"context"
	"encoding/json"
	"time"
)

type (
	// MsgHead 消息队列的头
	//todo 后续考虑用proto重构这个头
	MsgHead struct {
		Trace     string //追踪tid
		Timestamp int64  //发送时毫秒级时间戳
		Data      []byte //传送的内容
	}

	EventHandle interface {
		GetCtx() context.Context
		GetTs() time.Time
		GetData() []byte
	}
)

func NewEventMsg(ctx context.Context, data []byte) []byte {
	msg := MsgHead{
		Trace:     "",
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return msgBytes
}

func GetEventMsg(data []byte) EventHandle {
	msg := MsgHead{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil
	}
	return &msg
}

func (m *MsgHead) GetCtx() context.Context {
	//todo 等待实现链路追踪
	return context.Background()
}

func (m *MsgHead) GetTs() time.Time {
	return time.UnixMilli(m.Timestamp)
}

func (m *MsgHead) GetData() []byte {
	return m.Data
}
