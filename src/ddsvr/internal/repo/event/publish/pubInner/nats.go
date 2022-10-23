package pubInner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/nats-io/nats.go"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

func newNatsClient(conf conf.NatsConf) (PubInner, error) {
	nc, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) DevPubThing(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpThing, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubOta(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpOta, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubConfig(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpConfig, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubShadow(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpShadow, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubSDKLog(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpSDKLog, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error {
	str, _ := json.Marshal(info)
	switch conn {
	case Connect:
		return n.publish(ctx, topics.DeviceUpStatusConnected, str)
	case DisConnect:
		return n.publish(ctx, topics.DeviceUpStatusDisconnected, str)
	default:
		panic("not support conn type")
	}
}

func (n *NatsClient) publish(ctx context.Context, topic string, payload []byte) error {
	err := n.client.Publish(topic, events.NewEventMsg(ctx, payload))
	return err
}
