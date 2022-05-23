package innerLink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
)

//topic 定义
const (
	// TopicDevPubThing dd模块收到设备的发布消息后向内部推送以下topic 最后两个是产品id和设备名称
	TopicDevPubThing    = "dd.thing.device.clients.publish.thing.%s.%s"
	TopicDevPubThingAll = "dd.thing.device.clients.publish.thing.>"

	// TopicDevConnected dd模块收到设备的登录消息后向内部推送以下topic
	TopicDevConnected = "dd.thing.device.clients.connected"
	// TopicDevDisconnected dd模块收到设备的登出消息后向内部推送以下topic
	TopicDevDisconnected = "dd.thing.device.clients.disconnected"
	// TopicInnerPublish dd模块订阅以下topic,收到内部的发布消息后向设备推送
	TopicInnerPublish = "dd.thing.inner.publish"
)
const (
	ThingsDeliverGroup = "things_dm_group"
)

func NewNatsClient(conf conf.NatsConf) (*NatsClient, error) {
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

func (n *NatsClient) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	msg := events.NewEventMsg(ctx, devices.PublishToDev(topic, payload))
	err := n.client.Publish(TopicInnerPublish, msg)
	return err
}

func (n *NatsClient) SubscribeDevSync(ctx context.Context, topic string) (*SubDev, error) {
	subscription, err := n.client.SubscribeSync(topic)
	if err != nil {
		return nil, err
	}
	return NewSubDev(subscription), nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(TopicDevPubThingAll, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := device.GetDevPublish(ctx, msg)
			if err != nil {
				return err
			}
			err = handle(ctx).Thing(ele)
			return err
		}))
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(TopicDevConnected, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := device.GetDevConnMsg(ctx, msg)
			if err != nil {
				return err
			}
			return handle(ctx).Connected(ele)
		}))
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(TopicDevDisconnected, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := device.GetDevConnMsg(ctx, msg)
			if err != nil {
				return err
			}
			return handle(ctx).Disconnected(ele)
		}))
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsClient) ReqToDeviceSync(ctx context.Context, reqTopic, respTopic string, req *deviceSend.DeviceReq,
	productID, deviceName string) (*deviceSend.DeviceResp, error) {
	payload, _ := json.Marshal(req)
	err := n.PublishToDev(ctx, reqTopic, payload)
	if err != nil {
		return nil, err
	}
	handle, err := n.SubscribeDevSync(ctx, fmt.Sprintf(TopicDevPubThing, productID, deviceName))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := handle.UnSubscribe()
		if err != nil {
			logx.WithContext(ctx).Errorf("ReqToDeviceSync|UnSubscribe failure err:%v", err)
		}
	}()
	dead := utils.GetDeadLine(ctx, time.Now().Add(20*time.Second))
	for dead.After(time.Now()) {
		msg, err := handle.GetMsg(dead.Sub(time.Now()))
		if err != nil {
			return nil, err
		}
		if msg.Topic != respTopic { //不是订阅的topic
			continue
		}
		var dresp deviceSend.DeviceResp
		err = utils.Unmarshal(msg.Payload, &dresp)
		if err != nil { //如果是没法解析的说明不是需要的包,直接跳过即可
			continue
		}
		if dresp.ClientToken != req.ClientToken { //不是该请求的回复.跳过
			continue
		}
		return &dresp, nil
	}
	return nil, errors.DeviceTimeOut
}
