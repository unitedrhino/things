package subInner

import (
	"context"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/events"
	"gitee.com/i-Things/core/shared/events/topics"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

const (
	ThingsDDDeliverGroup = "things_dd_group"
)

func newNatsClient(conf conf.NatsConf) (SubInner, error) {
	nc, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) SubToDevMsg(handle Handle) error {
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
		}))
	return err
}
