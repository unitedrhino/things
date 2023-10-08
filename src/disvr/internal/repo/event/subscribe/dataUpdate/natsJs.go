package dataUpdate

import (
	"context"
	"encoding/json"
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
	natsJsConsumerName = "disvr"
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	nc, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: nc}, nil
}

func (n *NatsJsClient) Subscribe(handle Handle) error {
	_, err := n.client.Subscribe(topics.DmProductSchemaUpdate,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).ProductSchemaUpdate(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DmProductSchemaUpdate)))
	if err != nil {
		return err
	}
	_, err = n.client.Subscribe(topics.DmDeviceLogLevelUpdate,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceLogLevelUpdate(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DmDeviceLogLevelUpdate)))
	if err != nil {
		return err
	}
	_, err = n.client.Subscribe(topics.DmDeviceGatewayUpdate,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.GatewayUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceGatewayUpdate(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DmDeviceGatewayUpdate)))
	if err != nil {
		return err
	}

	_, err = n.client.Subscribe(topics.DmDeviceRemoteConfigUpdate,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceRemoteConfigUpdate(&tempInfo)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DmDeviceRemoteConfigUpdate)))
	if err != nil {
		return err
	}
	return nil
}
