package pubDev

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
)

type (
	MqttClient struct {
		client mqtt.Client
	}
)

func newEmqClient(conf *conf.MqttConf) (PubDev, error) {
	mc, err := clients.NewMqttClient(conf)
	if err != nil {
		return nil, err
	}
	return &MqttClient{
		client: mc,
	}, nil
}

func (d *MqttClient) Publish(ctx context.Context, topic string, payload []byte) error {
	return d.client.Publish(topic, 1, false, payload).Error()
}
