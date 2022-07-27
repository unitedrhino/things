package innerLink

import (
	//"github.com/i-Things/things/src/ddsvr/internal/domain"
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
)

type ConnType int8

const (
	Connect ConnType = iota
	DisConnect
)

type (
	InnerLink interface {
		// DevPubThing 向内部发布设备发布的物模型消息
		DevPubThing(ctx context.Context, publishMsg *devices.DevPublish) error
		DevPubOta(ctx context.Context, publishMsg *devices.DevPublish) error
		DevPubShadow(ctx context.Context, publishMsg *devices.DevPublish) error
		DevPubSDKLog(ctx context.Context, publishMsg *devices.DevPublish) error
		DevPubConfig(ctx context.Context, publishMsg *devices.DevPublish) error
		// PubConn 向内部发布连接及断连消息
		PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error
		Subscribe(handle Handle) error
	}
	Handle         func(ctx context.Context) InnerSubHandle
	InnerSubHandle interface {
		PublishToDev(topic string, payload []byte) error
	}
)

func NewInnerLink(c conf.InnerLinkConf) (InnerLink, error) {
	switch c.Mode {
	case conf.InnerLinkModeNats:
		return NewNatsClient(c.Nats)
	}
	return NewDirect()
}
