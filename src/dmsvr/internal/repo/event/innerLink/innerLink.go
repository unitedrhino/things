package innerLink

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceSend"
)

type (
	InnerLink interface {
		PublishToDev(ctx context.Context, topic string, payload []byte) error
		Subscribe(handle Handle) error
	}
	Handle         func(ctx context.Context) InnerSubHandle
	InnerSubHandle interface {
		Publish(out *deviceSend.Elements) error
		Connected(out *deviceSend.Elements) error
		Disconnected(out *deviceSend.Elements) error
	}
)

func NewInnerLink(conf config.InnerLinkConf) (InnerLink, error) {
	return NewNatsClient(conf.Nats)
}
