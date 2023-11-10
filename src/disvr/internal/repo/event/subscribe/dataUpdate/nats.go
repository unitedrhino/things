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
		client *clients.NatsClient
	}
)

var (
	natsJsConsumerName = "disvr"
)

func newNatsClient(conf conf.EventConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, natsJsConsumerName, conf.Nats)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	err := n.client.Subscribe(topics.DmProductSchemaUpdate,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).ProductSchemaUpdate(&tempInfo)
		})
	if err != nil {
		return err
	}
	err = n.client.Subscribe(topics.DmDeviceLogLevelUpdate,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceLogLevelUpdate(&tempInfo)
		})
	if err != nil {
		return err
	}
	err = n.client.Subscribe(topics.DmDeviceGatewayUpdate,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.GatewayUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceGatewayUpdate(&tempInfo)
		})
	if err != nil {
		return err
	}

	err = n.client.Subscribe(topics.DmDeviceRemoteConfigUpdate,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := events.DeviceUpdateInfo{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).DeviceRemoteConfigUpdate(&tempInfo)
		})
	if err != nil {
		return err
	}
	return nil
}
