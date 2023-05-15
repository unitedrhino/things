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
	_, err := n.client.Subscribe(topics.DmProductUpdateSchema,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).ProductSchemaUpdate(&tempInfo)
		}))
	if err != nil {
		return err
	}
	_, err = n.client.Subscribe(topics.DmDeviceUpdateLogLevel,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceLogLevelUpdate(&tempInfo)
		}))
	if err != nil {
		return err
	}
	_, err = n.client.Subscribe(topics.DmDeviceUpdateGateway,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.GatewayUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceGatewayUpdate(&tempInfo)
		}))
	if err != nil {
		return err
	}

	_, err = n.client.Subscribe(topics.DmDeviceUpdateRemoteConfig,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceRemoteConfigUpdate(&tempInfo)
		}))
	if err != nil {
		return err
	}
	return nil
}
