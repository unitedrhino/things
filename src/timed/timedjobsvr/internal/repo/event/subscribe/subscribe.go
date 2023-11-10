package subscribe

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
)

type (
	SubApp interface {
		Subscribe(handle Handle) error
	}
	Handle      func(ctx context.Context) ServerEvent
	ServerEvent interface {
		DataClean() error
	}
)

func NewSubServer(c conf.EventConf) (SubApp, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)

}
