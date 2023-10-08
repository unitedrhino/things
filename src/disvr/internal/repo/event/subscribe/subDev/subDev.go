package subDev

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceStatus"
)

type (
	SubDev interface {
		Subscribe(handle Handle) error
	}
	Handle        func(ctx context.Context) InnerSubEvent
	InnerSubEvent interface {
		// Thing 物模型消息
		Thing(out *deviceMsg.PublishMsg) error
		// SDKLog 设备调试日志
		SDKLog(out *deviceMsg.PublishMsg) error
		// Config 设备远程配置
		Config(out *deviceMsg.PublishMsg) error
		// Shadow 设备影子
		Shadow(out *deviceMsg.PublishMsg) error
		// Gateway 网关与子设备
		Gateway(out *deviceMsg.PublishMsg) error
		// Ota ota升级
		Ota(out *deviceMsg.PublishMsg) error
		// ext
		Ext(out *deviceMsg.PublishMsg) error

		Connected(out *deviceStatus.ConnectMsg) error
		Disconnected(out *deviceStatus.ConnectMsg) error
	}
)

func NewSubDev(c conf.EventConf) (SubDev, error) {
	switch c.Mode {
	case conf.EventModeNats:
		return newNatsClient(c.Nats)
	case conf.EventModeNatsJs:
		return newNatsJsClient(c.Nats)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
