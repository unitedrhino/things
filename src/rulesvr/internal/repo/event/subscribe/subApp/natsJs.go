package subApp

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
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
	_, err := n.client.QueueSubscribe(topics.ApplicationDeviceReportThingEventAllDevice, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			var stu application.EventReport
			err := utils.Unmarshal(msg, &stu)
			if err != nil {
				logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.[%s].Unmarshal err:%v", natsMsg.Subject, err)
			}
			return handle(ctx).DeviceEventReport(&stu)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.ApplicationDeviceReportThingEventAllDevice)))
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(topics.ApplicationDeviceReportThingPropertyAllDevice, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			var stu application.PropertyReport
			err := utils.Unmarshal(msg, &stu)
			if err != nil {
				logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.[%s].Unmarshal err:%v", natsMsg.Subject, err)
			}
			return handle(ctx).DevicePropertyReport(&stu)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.ApplicationDeviceReportThingPropertyAllDevice)))
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(topics.ApplicationDeviceStatusConnectedAllDevice, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			var stu application.ConnectMsg
			err := utils.Unmarshal(msg, &stu)
			if err != nil {
				logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.[%s].Unmarshal err:%v", natsMsg.Subject, err)
			}
			return handle(ctx).DeviceStatusConnected(&stu)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.ApplicationDeviceStatusConnectedAllDevice)))
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(topics.ApplicationDeviceStatusDisConnectedAllDevice, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			var stu application.ConnectMsg
			err := utils.Unmarshal(msg, &stu)
			if err != nil {
				logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.[%s].Unmarshal err:%v", natsMsg.Subject, err)
			}
			return handle(ctx).DeviceStatusDisConnected(&stu)
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.ApplicationDeviceStatusDisConnectedAllDevice)))
	if err != nil {
		return err
	}
	return nil
}
