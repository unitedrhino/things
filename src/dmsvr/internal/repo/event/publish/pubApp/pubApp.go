package pubApp

import (
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/domain/application"
	"gitee.com/i-Things/core/shared/errors"
)

type (
	PubApp interface {
		DeviceStatusConnected(ctx context.Context, msg application.ConnectMsg) error
		DeviceStatusDisConnected(ctx context.Context, msg application.ConnectMsg) error
		//应用事件通知-设备物模型属性上报通知
		DeviceThingPropertyReport(ctx context.Context, msg application.PropertyReport) error
		DeviceThingActionReport(ctx context.Context, msg application.ActionReport) error
		DeviceThingEventReport(ctx context.Context, msg application.EventReport) error
	}
)

func NewPubApp(c conf.EventConf) (PubApp, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	case conf.EventModeNatsJs:
		return newNatsJsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
