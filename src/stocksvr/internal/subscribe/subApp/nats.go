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
	_, err := n.client.QueueSubscribe(topics.ApplicationDeviceReportThingEventAllDevice, "",
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
	return nil
}
