package pubDev

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/things/share/clients"
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
	logx.WithContext(ctx).Infof("urPublish topic:%v payload: %s", topic, string(payload))
	err := d.client.Publish(topic, 1, false, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Publish failure err:%v topic:%v", err, topic)
	}
	return err
}
