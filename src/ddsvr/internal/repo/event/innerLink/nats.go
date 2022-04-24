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
		client nats.JetStreamContext
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
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name: ddExport.ThingsStreamName,
		Subjects: []string{
			ddExport.TopicInnerPublish,
			ddExport.TopicDevPublishAll,
			ddExport.TopicDevConnected,
			ddExport.TopicDevDisconnected,
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = js.AddConsumer(ddExport.ThingsStreamName, &nats.ConsumerConfig{
		Durable:   ddExport.ThingsConsumeName,
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: js}, nil
}

func (n *NatsClient) Publish(ctx context.Context, topic string, payload []byte) error {

	_, err := n.client.Publish(topic, events.NewEventMsg(ctx, payload))
	return err
}
func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(ddExport.TopicInnerPublish, ddExport.SvrName, func(msg *nats.Msg) {
		msg.Ack()
		ctx, topic, payload := ddExport.GetPublish(msg.Data)
		handle(ctx).PublishToDev(topic, payload)
	})
	return err
}
