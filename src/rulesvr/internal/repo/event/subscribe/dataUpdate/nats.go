package dataUpdate

import (
	"context"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/events"
	"gitee.com/i-Things/core/shared/events/topics"
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

func (n *NatsClient) Subscribe(handle Handle) error {
	if _, err := n.client.Subscribe(topics.DmProductSchemaUpdate,
		events.NatsSubWithType(func(ctx context.Context, tempInfo events.DeviceUpdateInfo, natsMsg *nats.Msg) error {
			return handle(ctx).ProductSchemaUpdate(&tempInfo)
		})); err != nil {
		return err
	}
	if _, err := n.client.Subscribe(topics.RuleSceneInfoUpdate,
		events.NatsSubWithType(func(ctx context.Context, tempInfo events.ChangeInfo, natsMsg *nats.Msg) error {
			return handle(ctx).SceneInfoUpdate(&tempInfo)
		})); err != nil {
		return err
	}
	if _, err := n.client.Subscribe(topics.RuleSceneInfoDelete,
		events.NatsSubWithType(func(ctx context.Context, tempInfo events.ChangeInfo, natsMsg *nats.Msg) error {
			return handle(ctx).SceneInfoDelete(&tempInfo)
		})); err != nil {
		return err
	}
	return nil
}
