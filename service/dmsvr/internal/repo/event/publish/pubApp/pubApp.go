package pubApp

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/domain/application"
	"gitee.com/unitedrhino/share/errors"
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
