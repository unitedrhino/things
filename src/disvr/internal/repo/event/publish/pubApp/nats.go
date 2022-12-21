package pubApp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
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

func (n *NatsClient) DeviceThingPropertyReport(ctx context.Context, msg application.PropertyReport) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingProperty, msg.Device.ProductID, msg.Device.DeviceName, msg.Identifier)
	err := n.client.Publish(topic, pubMsg)
	return err
}

func (n *NatsClient) DeviceThingEventReport(ctx context.Context, msg application.EventReport) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingEvent,
		msg.Device.ProductID, msg.Device.DeviceName, msg.Type, msg.Identifier)
	err := n.client.Publish(topic, pubMsg)
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