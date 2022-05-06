package innerLink

import (
	//"github.com/i-Things/things/src/ddsvr/internal/domain"
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/ddsvr/internal/config"
)

type ConnType int8

const (
	Connect ConnType = iota
	DisConnect
)

type (
	InnerLink interface {
		//向内部发布设备发布的消息
		PubDevPublish(ctx context.Context, publishMsg devices.DevPublish) error
		//向内部发布连接及断连消息
		PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error
		Subscribe(handle Handle) error
	}
	Handle         func(ctx context.Context) InnerSubHandle
	InnerSubHandle interface {
		PublishToDev(topic string, payload []byte) error
	}
)

func NewInnerLink(conf config.InnerLinkConf) (InnerLink, error) {
	return NewNatsClient(conf.Nats)
}
