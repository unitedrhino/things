package dataUpdate

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/nats-io/nats.go"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

var (
	natsJsConsumerName = "rulesvr"
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	nc, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: nc}, nil
}

func (n *NatsJsClient) Subscribe(handle Handle) error {
	if _, err := n.client.Subscribe(topics.DmProductSchemaUpdate,
		events.NatsSubWithType(func(ctx context.Context, tempInfo events.DeviceUpdateInfo, natsMsg *nats.Msg) error {
			return handle(ctx).ProductSchemaUpdate(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DmProductSchemaUpdate))); err != nil {
		return err
	}
	if _, err := n.client.Subscribe(topics.RuleSceneInfoUpdate,
		events.NatsSubWithType(func(ctx context.Context, tempInfo events.ChangeInfo, natsMsg *nats.Msg) error {
			return handle(ctx).SceneInfoUpdate(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.RuleSceneInfoUpdate))); err != nil {
		return err
	}
	if _, err := n.client.Subscribe(topics.RuleSceneInfoDelete,
		events.NatsSubWithType(func(ctx context.Context, tempInfo events.ChangeInfo, natsMsg *nats.Msg) error {
			return handle(ctx).SceneInfoDelete(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.RuleSceneInfoDelete))); err != nil {
		return err
	}
	return nil
}
