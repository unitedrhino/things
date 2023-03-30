package dataUpdate

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
)

type (
	DataUpdate interface {
		Subscribe(handle Handle) error
	}
	Handle       func(ctx context.Context) UpdateHandle
	UpdateHandle interface {
		ProductSchemaUpdate(info *events.DataUpdateInfo) error
	}
)

func NewDataUpdate(c conf.EventConf) (DataUpdate, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
