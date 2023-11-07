package server

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
)

type (
	Server interface {
		Subscribe(handle Handle) error
	}
	Handle       func(ctx context.Context) ServerHandle
	ServerHandle interface {
		ActionCheck(req *deviceMsg.PublishMsg) error
	}
)

func NewServer(c conf.EventConf) (Server, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
