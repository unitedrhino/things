package subDev

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
)

type (
	SubDev interface {
		SubDevMsg(handle Handle) error
	}
	Handle       func(ctx context.Context) DevSubHandle
	DevSubHandle interface {
		// Ota ota远程升级
		Ota(topic string, payload []byte) error
		// Config 设备远程配置
		Config(topic string, payload []byte) error
		// Shadow 设备影子
		Shadow(topic string, payload []byte) error
		// Log 设备调试日志
		SDKLog(topic string, payload []byte) error
		// Gateway 网关子设备
		Gateway(topic string, payload []byte) error
		// Thing 物模型消息
		Thing(topic string, payload []byte) error
		Connected(out *devices.DevConn) error
		Disconnected(out *devices.DevConn) error
	}
)

func Check(conf conf.DevLinkConf) error {
	if conf.Mqtt == nil {
		return fmt.Errorf("DevLinkConf need")
	}
	return nil
}

func NewSubDev(conf conf.DevLinkConf) (SubDev, error) {
	if err := Check(conf); err != nil {
		return nil, err
	}
	return newEmqClient(conf.Mqtt)
}
