package pubApp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/src/disvr/internal/domain/service/application"
	"github.com/nats-io/nats.go"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

func newNatsClient(conf conf.NatsConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	msg := events.NewEventMsg(ctx, devices.PublishToDev(topic, payload))
	err := n.client.Publish(topics.DeviceDownAll, msg)
	return err
}

func (n *NatsClient) DeviceStatusConnected(ctx context.Context, msg application.ConnectMsg) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusConnected, msg.Device.ProductID, msg.Device.DeviceName)
	err := n.client.Publish(topic, pubMsg)
	return err
}

func (n *NatsClient) DeviceStatusDisConnected(ctx context.Context, msg application.ConnectMsg) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusDisConnected, msg.Device.ProductID, msg.Device.DeviceName)
	err := n.client.Publish(topic, pubMsg)
	return err
}
