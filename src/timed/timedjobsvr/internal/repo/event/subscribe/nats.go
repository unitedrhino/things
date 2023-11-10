package subscribe

import (
	"context"
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

const (
	ThingsDeliverGroup = "things_timed_job_group"
	natsJsConsumerName = "timedjobsvr"
)

var ()

func newNatsClient(conf conf.EventConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, natsJsConsumerName, conf.Nats)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	err := n.client.QueueSubscribe(topics.TimedJobClean, ThingsDeliverGroup,
		func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			return handle(ctx).DataClean()
		})
	if err != nil {
		return err
	}
	return nil
}
