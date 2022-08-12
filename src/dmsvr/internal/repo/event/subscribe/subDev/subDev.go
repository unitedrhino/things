package subDev

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
)

type (
	SubDev interface {
		Subscribe(handle Handle) error
	}
	Handle        func(ctx context.Context) InnerSubEvent
	InnerSubEvent interface {
		// Thing 物模型消息
		Thing(out *device.PublishMsg) error
		// SDK Log 设备调试日志
		SDKLog(out *device.PublishMsg) error
		// Config 设备远程配置
		Config(out *device.PublishMsg) error
		// Shadow 设备影子
		Shadow(out *device.PublishMsg) error
		// Ota ota升级
		Ota(out *device.PublishMsg) error
		Connected(out *device.ConnectMsg) error
		Disconnected(out *device.ConnectMsg) error
	}
)

func NewSubDev(c conf.EventConf) (SubDev, error) {
	if c.Mode == conf.EventModeNats {
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
