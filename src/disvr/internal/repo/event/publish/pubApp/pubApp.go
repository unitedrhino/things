package pubApp

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/errors"
)

type (
	PubApp interface {
		DeviceStatusConnected(ctx context.Context, msg application.ConnectMsg) error
		DeviceStatusDisConnected(ctx context.Context, msg application.ConnectMsg) error
		DeviceThingPropertyReport(ctx context.Context, msg application.PropertyReport) error
		DeviceThingEventReport(ctx context.Context, msg application.EventReport) error
	}
)

func NewPubApp(c conf.EventConf) (PubApp, error) {
	if c.Mode == conf.EventModeNats {
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
