package subInner

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

var (
	natsJsConsumerName = "ddsvr"
)

func newNatsJsClient(conf conf.NatsConf) (SubInner, error) {
	js, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: js}, nil
}

func (n *NatsJsClient) SubToDevMsg(handle Handle) error {
	_, err := n.client.QueueSubscribe(topics.DeviceDownAll, ThingsDDDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			//给设备回包之前，将链路信息span推送至jaeger
			_, span := ctxs.StartSpan(ctx, topics.DeviceDownAll, "")
			info := devices.GetPublish(msg)
			logx.WithContext(ctx).Infof("ddsvr.mqtt.SubDevMsg Handle:%s Type:%v Payload:%v",
				info.Handle, info.Type, string(info.Payload))
			defer span.End()
			err := handle(ctx).PublishToDev(info)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.PublishToDev failure err:%v", err)
			}
			return err
		}), nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topics.DeviceDownAll)))
	return err
}
