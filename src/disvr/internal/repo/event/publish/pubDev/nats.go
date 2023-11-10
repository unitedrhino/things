package pubDev

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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

func (n *NatsClient) PublishToDev(ctx context.Context, respMsg *deviceMsg.PublishMsg) error {
	msg := events.NewEventMsg(ctx, devices.PublishToDev(respMsg.Handle, respMsg.Type, respMsg.Payload, respMsg.ProductID, respMsg.DeviceName))
	logx.WithContext(ctx).Infof("PublishToDev sendMsg:%s", string(msg))
	err := n.client.Publish(fmt.Sprintf(topics.DeviceDownMsg, respMsg.Handle, respMsg.ProductID, respMsg.DeviceName), msg)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s Publish failure err:%v", utils.FuncName(), err)
	}
	return err
}

func (n *NatsClient) ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, compareMsg CompareMsg) ([]byte, error) {
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

func (n *NatsClient) subscribeDevSync(ctx context.Context, topic string) (*natsSubDev, error) {
	subscription, err := n.client.SubscribeSync(topic)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.SubscribeSync failure err:%v", utils.FuncName(), err)
		return nil, err
	}
	return newNatsSubDev(subscription), nil
}
