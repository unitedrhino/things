package pubDev

import (
	"context"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	MqttClient struct {
		client *clients.MqttClient
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
	err := d.client.Publish(topic, 1, false, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Publish failure err:%v topic:%v", err, topic)
	}
	return err
}
