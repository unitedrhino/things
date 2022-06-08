package devLink

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/ddsvr/internal/config"
)

type DevClient struct {
	client mqtt.Client
}

type (
	DevLink interface {
		Publish(ctx context.Context, topic string, payload []byte) error
		SubScribe(handle Handle) error
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
		// Thing 物模型消息
		Thing(topic string, payload []byte) error
		Connected(out *devices.DevConn) error
		Disconnected(out *devices.DevConn) error
	}
)

func Check(conf config.DevLinkConf) error {
	if conf.Mqtt == nil {
		return fmt.Errorf("DevLinkConf need")
	}
	return nil
}

func NewDevClient(conf config.DevLinkConf) (DevLink, error) {
	if err := Check(conf); err != nil {
		return nil, err
	}
	return NewEmqClient(conf.Mqtt)
}
