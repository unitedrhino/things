package subInner

import (
	"context"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/events/topics"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

var (
	natsJsConsumerName = "dgsvr"
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
			logx.WithContext(ctx).Infof("dgsvr.mqtt.SubDevMsg Handle:%s Type:%v Payload:%v",
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
