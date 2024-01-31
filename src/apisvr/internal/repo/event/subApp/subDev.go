package subApp

import (
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/domain/application"
	"gitee.com/i-Things/core/shared/errors"
)

type (
	SubApp interface {
		Subscribe(handle Handle) error
	}
	Handle      func(ctx context.Context) AppSubEvent
	AppSubEvent interface {
		// DeviceEventReport 设备事件上报
		DeviceEventReport(out *application.EventReport) error
		// DevicePropertyReport 设备属性上报
		DevicePropertyReport(out *application.PropertyReport) error
		// DeviceStatusConnected 设备连接状态上报
		DeviceStatusConnected(out *application.ConnectMsg) error
		// DeviceStatusDisConnected 设备离线上报
		DeviceStatusDisConnected(out *application.ConnectMsg) error
	}
)

func NewSubApp(c conf.EventConf) (SubApp, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	case conf.EventModeNatsJs:
		return newNatsJsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
