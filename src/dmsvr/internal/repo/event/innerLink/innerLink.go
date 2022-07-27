package innerLink

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	deviceSend "github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
)

type (
	InnerLink interface {
		PublishToDev(ctx context.Context, topic string, payload []byte) error
		Subscribe(handle Handle) error
		ReqToDeviceSync(ctx context.Context, reqTopic, respTopic string, req *deviceSend.DeviceReq,
			productID, deviceName string) (*deviceSend.DeviceResp, error)
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

func NewInnerLink(c conf.InnerLinkConf) (InnerLink, error) {
	if c.Mode == conf.InnerLinkModeNats {
		return NewNatsClient(c.Nats)
	}
	//todo 等待支持直接调用模式
	return nil, nil
}
