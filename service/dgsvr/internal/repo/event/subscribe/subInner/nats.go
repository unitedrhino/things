package subInner

import (
	"context"
	"fmt"
	"time"

	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/events/topics"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/protocols"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client *clients.NatsClient
	}
)

const (
	ThingsDDDeliverGroup = "things_dg_group"
)

func newNatsClient(conf conf.EventConf, nodeID int64) (SubInner, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, conf.Nats.Consumer, conf.Nats, nodeID)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) SubToDevMsg(handle Handle) error {
	topic := fmt.Sprintf(topics.DeviceDownAll, protocols.ProtocolCodeUrMqtt)
	_, err := n.client.QueueSubscribe(topic, ThingsDDDeliverGroup,
		func(ctx context.Context, ts time.Time, msg []byte, natsMsg *nats.Msg) error {
			//给设备回包之前，将链路信息span推送至jaeger
			_, span := ctxs.StartSpan(ctx, topic, "")
			info := devices.GetPublish(msg)
			logx.WithContext(ctx).Infof("dgsvr.mqtt.SubDevMsg Handle:%s Type:%v Payload:%v",
				info.Handle, info.Type, string(info.Payload))
			defer span.End()
			err := handle(ctx).PublishToDev(info)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.PublishToDev failure err:%v", err)
			}
			return err
		})
	return err
}
