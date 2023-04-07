package dataUpdate

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
)

type (
	DataUpdate interface {
		ProductScriptUpdate(ctx context.Context, info *events.DataUpdateInfo) error
		ProductSchemaUpdate(ctx context.Context, info *events.DataUpdateInfo) error
		DeviceLogLevelUpdate(ctx context.Context, info *events.DataUpdateInfo) error
		DeviceGatewayUpdate(ctx context.Context, info *events.GatewayUpdateInfo) error
		DeviceRemoteConfigUpdate(ctx context.Context, info *events.DataUpdateInfo) error
	}
)

func NewDataUpdate(c conf.EventConf) (DataUpdate, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
