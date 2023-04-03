package subInner

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
)

type (
	SubInner interface {
		SubToDevMsg(handle Handle) error
	}
	Handle         func(ctx context.Context) InnerSubHandle
	InnerSubHandle interface {
		PublishToDev(info *devices.InnerPublish) error
	}
)

func NewSubInner(conf conf.EventConf) (SubInner, error) {
	return newNatsClient(conf.Nats)
}
