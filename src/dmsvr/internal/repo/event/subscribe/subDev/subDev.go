package subDev

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
)

type (
	SubDev interface {
		Subscribe(handle Handle) error
	}
	Handle        func(ctx context.Context) InnerSubEvent
	InnerSubEvent interface {
		// Thing 物模型消息
		Thing(out *deviceMsg.PublishMsg) error
		// SDK Log 设备调试日志
		SDKLog(out *deviceMsg.PublishMsg) error
		// Config 设备远程配置
		Config(out *deviceMsg.PublishMsg) error
		// Shadow 设备影子
		Shadow(out *deviceMsg.PublishMsg) error
		// Ota ota升级
		Ota(out *deviceMsg.PublishMsg) error
		Connected(out *deviceMsg.ConnectMsg) error
		Disconnected(out *deviceMsg.ConnectMsg) error
	}
)

func NewSubDev(c conf.EventConf) (SubDev, error) {
	if c.Mode == conf.EventModeNats {
		return newNatsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
