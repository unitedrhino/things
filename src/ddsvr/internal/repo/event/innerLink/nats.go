package innerLink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/traces"
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

func (n *NatsClient) DevPubThing(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpThing, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubOta(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpOta, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubConfig(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpConfig, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubShadow(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpShadow, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) DevPubLog(ctx context.Context, publishMsg *devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(topics.DeviceUpLog, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error {
	str, _ := json.Marshal(info)
	switch conn {
	case Connect:
		return n.publish(ctx, topics.DeviceUpStatusConnected, str)
	case DisConnect:
		return n.publish(ctx, topics.DeviceUpStatusDisconnected, str)
	default:
		panic("not support conn type")
	}
}

func (n *NatsClient) publish(ctx context.Context, topic string, payload []byte) error {
	err := n.client.Publish(topic, events.NewEventMsg(ctx, payload))
	return err
}
func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(topics.DeviceDownAll, ThingsDDDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			topic, payload := devices.GetPublish(msg)
			//给设备回包之前，将链路信息span推送至jaeger
			_, span := traces.StartSpan(ctx, topics.DeviceDownAll, "")
			logx.Infof("[ddsvr.mqtt.SubScribe]|-------------------trace:%s, spanid:%s|topic:%s",
				span.SpanContext().TraceID(), span.SpanContext().SpanID(), topics.DeviceDownAll)
			defer span.End()
			return handle(ctx).PublishToDev(topic, payload)
		}))
	return err
}
