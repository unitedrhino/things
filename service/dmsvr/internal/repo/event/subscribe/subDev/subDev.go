package subDev

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceStatus"
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

func NewSubDev(c conf.EventConf, nodeID int64) (SubDev, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c, nodeID)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
}
