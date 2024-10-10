package server

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/share/events/topics"
	"github.com/nats-io/nats.go"
)

type (
	NatsClient struct {
		client *clients.NatsClient
	}
)

var (
	natsJsConsumerName = "dmsvr"
)

func newNatsClient(conf conf.EventConf, nodeID int64) (*NatsClient, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, natsJsConsumerName, conf.Nats, nodeID)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(topics.DmActionCheckDelay, natsJsConsumerName,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			tempInfo := deviceMsg.PublishMsg{}
			err := json.Unmarshal(msg, &tempInfo)
			if err != nil {
				return err
			}
			return handle(ctx).ActionCheck(&tempInfo)
		})
	if err != nil {
		return err
	}
	return nil
}
