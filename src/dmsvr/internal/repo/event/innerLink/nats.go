package innerLink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
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
	err := n.client.Publish(topics.DeviceDownAll, msg)
	return err
}

func (n *NatsClient) SubscribeDevSync(ctx context.Context, topic string) (*SubDev, error) {
	subscription, err := n.client.SubscribeSync(topic)
	if err != nil {
		return nil, err
	}
	return NewSubDev(subscription), nil
}

func (n *NatsClient) QueueSubscribeDevPublish(topic string,
	handleFunc func(ctx context.Context, msg *device.PublishMsg) error) error {
	_, err := n.client.QueueSubscribe(topic, ThingsDeliverGroup,
		events.NatsSubscription(func(ctx context.Context, msg []byte) error {
			ele, err := device.GetDevPublish(ctx, msg)
			if err != nil {
				return err
			}
			return handleFunc(ctx, ele)
		}))
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	err := n.QueueSubscribeDevPublish(topics.DeviceUpThingAll, func(ctx context.Context, msg *device.PublishMsg) error {
		err := handle(ctx).Thing(msg)
		return err
	})
	if err != nil {
		return err
	}
	err = n.QueueSubscribeDevPublish(topics.DeviceUpOtaAll, func(ctx context.Context, msg *device.PublishMsg) error {
		err := handle(ctx).Ota(msg)
		return err
	})
	if err != nil {
		return err
	}
	err = n.QueueSubscribeDevPublish(topics.DeviceUpConfigAll, func(ctx context.Context, msg *device.PublishMsg) error {
		err := handle(ctx).Config(msg)
		return err
	})
	if err != nil {
		return err
	}
	err = n.QueueSubscribeDevPublish(topics.DeviceUpSDKLogAll, func(ctx context.Context, msg *device.PublishMsg) error {
		err := handle(ctx).SDKLog(msg)
		return err
	})
	if err != nil {
		return err
	}
	err = n.QueueSubscribeDevPublish(topics.DeviceUpShadowAll, func(ctx context.Context, msg *device.PublishMsg) error {
		err := handle(ctx).Shadow(msg)
		return err
	})
	if err != nil {
		return err
	}

	_, err = n.client.QueueSubscribe(topics.DeviceUpStatusConnected, ThingsDeliverGroup,
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

	_, err = n.client.QueueSubscribe(topics.DeviceUpStatusDisconnected, ThingsDeliverGroup,
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
	handle, err := n.SubscribeDevSync(ctx, fmt.Sprintf(topics.DeviceUpThing, productID, deviceName))
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
