package subApp

import (
	"context"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/domain/application"
	"gitee.com/i-Things/core/shared/events"
	"gitee.com/i-Things/core/shared/events/topics"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

const (
	ThingsDeliverGroup = "things_api_group"
)

func newNatsClient(conf conf.NatsConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(topics.ApplicationDeviceReportThingEventAllDevice, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			var stu application.EventReport
			err := utils.Unmarshal(msg, &stu)
			if err != nil {
				logx.WithContext(ctx).Errorf("Subscribe.QueueSubscribe.[%s].Unmarshal err:%v", natsMsg.Subject, err)
			}
			return handle(ctx).DeviceEventReport(&stu)
		}))
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
		}))
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
		}))
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
		}))
	if err != nil {
		return err
	}
	return nil
}
