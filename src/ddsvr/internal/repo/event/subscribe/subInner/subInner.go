package subInner

import (
	"context"
	"github.com/i-Things/things/shared/conf"
)

type (
	SubInner interface {
		SubToDevMsg(handle Handle) error
	}
	Handle         func(ctx context.Context) InnerSubHandle
	InnerSubHandle interface {
		PublishToDev(topic string, payload []byte) error
	}
)

func NewSubInner(conf conf.EventConf) (SubInner, error) {
	return newNatsClient(conf.Nats)
}
