package dataUpdate

import (
	"context"
	"github.com/i-Things/things/src/dmsvr/internal/config"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
)

type (
	DataUpdate interface {
		TempModelUpdate(ctx context.Context, info *thing.TemplateInfo) error
		Subscribe(handle Handle) error
	}
	Handle             func(ctx context.Context) DataUpdateSubEvent
	DataUpdateSubEvent interface {
		TempModelClearCache(info *thing.TemplateInfo) error
	}
)

func NewDataUpdate(conf config.InnerLinkConf) (DataUpdate, error) {
	return NewNatsClient(conf.Nats)
}
