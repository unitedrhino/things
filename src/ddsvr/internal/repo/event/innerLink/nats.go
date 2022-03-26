package innerLink

import (
	"context"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/nats-io/nats.go"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

func NewNatsClient(conf conf.NatsConf) (InnerLink, error) {
	connectOpts := nats.Options{
		Url:      conf.Url,
		User:     conf.User,
		Password: conf.Pass,
		Token:    conf.Token,
	}
	nc, err := connectOpts.Connect()
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Publish(ctx context.Context, topic string, payload []byte) error {
	return n.client.Publish(topic, events.NewEventMsg(ctx, payload))
}
func (n *NatsClient) Subscribe(handle Handle) error {
	n.client.QueueSubscribe(ddExport.TopicInnerPublish, ddExport.SvrName, func(msg *nats.Msg) {
		ctx, topic, payload := ddExport.GetPublish(msg.Data)
		handle(ctx).PublishToDev(topic, payload)
	})
	return nil
}
