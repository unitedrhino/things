package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/nats-io/nats.go"
)

type (
	NatsClient struct {
		client *clients.NatsClient
	}
)

var (
	natsJsConsumerName = "vidsvr"
)

func newNatsClient(conf conf.EventConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, natsJsConsumerName, conf.Nats)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	err := n.client.QueueSubscribe(topics.VidInfoCheckStatus, natsJsConsumerName,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			jsonStr, _ := json.Marshal(natsMsg)
			fmt.Println("[******]   QueueSubscribe1", "data:", string(jsonStr))
			return handle(ctx).ActionCheck()
		})
	if err != nil {
		return err
	}
	err = n.client.QueueSubscribe(topics.VidInfoInitDatabase, natsJsConsumerName,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			jsonStr, _ := json.Marshal(natsMsg)
			fmt.Println("[***Once***]   QueueSubscribe2 ", "data:", string(jsonStr))
			return handle(ctx).ActionCheck()
		})
	if err != nil {
		return err
	}
	return nil
}
