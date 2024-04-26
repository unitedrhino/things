package pubDev

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/events/topics"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	NatsClient struct {
		client *clients.NatsClient
	}
)

func newNatsClient(conf conf.EventConf) (*NatsClient, error) {
	nc, err := clients.NewNatsClient2(conf.Mode, conf.Nats.Consumer, conf.Nats)
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: nc}, nil
}

func (n *NatsClient) PublishToDev(ctx context.Context, respMsg *deviceMsg.PublishMsg) error {
	startTime := time.Now()
	if respMsg.ProtocolCode == "" {
		respMsg.ProtocolCode = def.ProtocolCodeIThings
	}
	defer func() {
		logx.WithContext(ctx).WithDuration(time.Now().Sub(startTime)).
			Infof("PublishToDev startTime:%v", startTime)
	}()
	err := n.client.Publish(ctx, fmt.Sprintf(topics.DeviceDownMsg, respMsg.ProtocolCode, respMsg.Handle, respMsg.ProductID, respMsg.DeviceName), devices.PublishToDev(
		respMsg.Handle, respMsg.Type, respMsg.Payload, respMsg.ProtocolCode,
		respMsg.ProductID, respMsg.DeviceName))

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
	handle, err := n.subscribeDevSync(ctx, fmt.Sprintf(topics.DeviceUpMsg, reqMsg.Type, reqMsg.ProductID, reqMsg.DeviceName))
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
