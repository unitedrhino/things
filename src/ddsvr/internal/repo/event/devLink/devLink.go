package devLink

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/src/ddsvr/ddDef"
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
		Publish(topic string, payload []byte) error
		Login(out *ddDef.DevLogInOut) error
		Logout(out *ddDef.DevLogInOut) error
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
