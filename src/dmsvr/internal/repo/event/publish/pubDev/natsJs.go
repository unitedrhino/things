package pubDev

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/events"
	"gitee.com/i-Things/core/shared/events/topics"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/google/uuid"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

var (
	natsJsConsumerName = "disvr"
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	nc, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &NatsJsClient{client: nc}, nil
}

func (n *NatsJsClient) PublishToDev(ctx context.Context, respMsg *deviceMsg.PublishMsg) error {
	startTime := time.Now()
	msg := events.NewEventMsg(ctx, devices.PublishToDev(respMsg.Handle, respMsg.Type, respMsg.Payload, respMsg.ProductID, respMsg.DeviceName))
	defer func() {
		logx.WithContext(ctx).WithDuration(time.Now().Sub(startTime)).Infof("PublishToDev startTime:%v sendMsg:%s", startTime, string(msg))
	}()
	_, err := n.client.Publish(fmt.Sprintf(topics.DeviceDownMsg, respMsg.Handle, respMsg.ProductID, respMsg.DeviceName), msg)

	if err != nil {
		logx.WithContext(ctx).Errorf("%s Publish failure err:%v", utils.FuncName(), err)
	}
	return err
}

func (n *NatsJsClient) ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, compareMsg CompareMsg) ([]byte, error) {
	err := n.PublishToDev(ctx, reqMsg)
	if err != nil {
		return nil, err
	}
	handle, err := n.subscribeDevSync(ctx, fmt.Sprintf(topics.DeviceUpThing, reqMsg.ProductID, reqMsg.DeviceName))
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.subscribeDevSync failure err:%v", utils.FuncName(), err)
		return nil, err
	}
	defer func() {
		err := handle.UnSubscribe()
		if err != nil {
			logx.WithContext(ctx).Errorf("ReqToDeviceSync.UnSubscribe failure err:%v", err)
		}
	}()
	dead := time.Now().Add(10 * time.Second)
	for dead.After(time.Now()) {
		msg, err := handle.GetMsg(dead.Sub(time.Now()))
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.GetMsg failure err:%v", utils.FuncName(), err)
			return nil, err
		}
		if msg.Handle != reqMsg.Handle || msg.Type != reqMsg.Type { //不是订阅的topic
			continue
		}
		if !compareMsg(msg.Payload) {
			continue
		}
		return msg.Payload, nil
	}
	return nil, errors.DeviceTimeOut
}

func (n *NatsJsClient) subscribeDevSync(ctx context.Context, topic string) (*natsSubDev, error) {
	subscription, err := n.client.SubscribeSync(topic, nats.Durable(events.GenNatsJsDurable(natsJsConsumerName, topic+"--"+uuid.NewString())))
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.SubscribeSync failure err:%v", utils.FuncName(), err)
		return nil, err
	}
	return newNatsSubDev(subscription), nil
}
