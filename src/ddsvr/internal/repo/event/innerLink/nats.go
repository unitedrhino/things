package innerLink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/traces"
	"github.com/i-Things/things/src/ddsvr/dd"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

const (
	ThingsQueueConsumeName = "things_dd_queue_consume"
	//topic 定义
	ThingsStreamName = "thing_msg"
	// TopicDevPublish dd模块收到设备的发布消息后向内部推送以下topic 最后两个是产品id和设备名称
	TopicDevPublish = "dd.thing.device.clients.publish.%s.%s"

	// TopicDevConnected dd模块收到设备的登录消息后向内部推送以下topic
	TopicDevConnected = "dd.thing.device.clients.connected"
	// TopicDevDisconnected dd模块收到设备的登出消息后向内部推送以下topic
	TopicDevDisconnected = "dd.thing.device.clients.disconnected"
	// TopicInnerPublish dd模块订阅以下topic,收到内部的发布消息后向设备推送
	TopicInnerPublish = "dd.thing.inner.publish"
	TopicThing        = "dd.thing.device.clients.>"
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

func (n *NatsClient) PubDevPublish(ctx context.Context, publishMsg devices.DevPublish) error {
	pubStr, _ := json.Marshal(publishMsg)
	return n.publish(ctx,
		fmt.Sprintf(TopicDevPublish, publishMsg.ProductID, publishMsg.DeviceName), pubStr)
}

func (n *NatsClient) PubConn(ctx context.Context, conn ConnType, info *devices.DevConn) error {
	str, _ := json.Marshal(info)
	switch conn {
	case Connect:
		return n.publish(ctx, TopicDevConnected, str)
	case DisConnect:
		return n.publish(ctx, TopicDevDisconnected, str)
	default:
		panic("not support conn type")
	}
}

func (n *NatsClient) publish(ctx context.Context, topic string, payload []byte) error {
	err := n.client.Publish(topic, events.NewEventMsg(ctx, payload))
	return err
}
func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(TopicInnerPublish, dd.ThingsDDDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			topic, payload := devices.GetPublish(msg)
			//给设备回包之前，将链路信息span推送至jaeger
			_, span := traces.StartSpan(ctx, TopicInnerPublish, "")
			logx.Infof("[mqtt.SubScribe]|-------------------trace:%s, spanid:%s|topic:%s",
				span.SpanContext().TraceID(), span.SpanContext().SpanID(), TopicInnerPublish)
			defer span.End()
			return handle(ctx).PublishToDev(topic, payload)
		}))
	return err
}
