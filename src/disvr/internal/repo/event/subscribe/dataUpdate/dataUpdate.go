package dataUpdate

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
)

type (
	DataUpdate interface {
		Subscribe(handle Handle) error
	}
	Handle           func(ctx context.Context) DataUpdateHandle
	DataUpdateHandle interface {
		SchemaClearCache(info *schema.Info) error
	}
)

func NewDataUpdate(c conf.EventConf) (DataUpdate, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
