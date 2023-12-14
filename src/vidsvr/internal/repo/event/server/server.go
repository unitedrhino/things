package server

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
)

type (
	Server interface {
		Subscribe(handle Handle) error
	}
	Handle       func(ctx context.Context) ServerHandle
	ServerHandle interface {
		ActionCheck() error
		ActionInit() error
	}
)

func NewServer(c conf.EventConf) (Server, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
