package dataUpdate

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
)

type (
	DataUpdate interface {
		TempModelUpdate(ctx context.Context, info *schema.SchemaInfo) error
		Subscribe(handle Handle) error
	}
	Handle             func(ctx context.Context) DataUpdateSubEvent
	DataUpdateSubEvent interface {
		TempModelClearCache(info *schema.SchemaInfo) error
	}
)

func NewDataUpdate(c conf.InnerLinkConf) (DataUpdate, error) {
	switch c.Mode {
	case conf.InnerLinkModeNats:
		return NewNatsClient(c.Nats)
	}
	return NewDirect()
}
