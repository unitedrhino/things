package pubInner

import (
	//"github.com/i-Things/things/src/dgsvr/internal/domain"
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/errors"
)

type ConnType int8

const (
	Connect ConnType = iota
	DisConnect
)

type (
	PubInner interface {
		DevPubMsg(ctx context.Context, publishMsg *devices.DevPublish) error
		// PubConn 向内部发布连接及断连消息
		PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error
	}
)

func NewPubInner(c conf.EventConf) (PubInner, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	case conf.EventModeNatsJs:
		return newNatsJsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
