package pubDev

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/events/topics"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/protocols"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	pubDevClient struct {
		client *eventBus.FastEvent
	}
)

func newPubDevClient(fast *eventBus.FastEvent) (*pubDevClient, error) {
	return &pubDevClient{client: fast}, nil
}

func (n *pubDevClient) PublishToDev(ctx context.Context, reqMsg *deviceMsg.PublishMsg) error {
	startTime := time.Now()
	if reqMsg.ProtocolCode == "" {
		reqMsg.ProtocolCode = protocols.ProtocolCodeUrMqtt
	}
	reqMsg = s.DownBeforeTrans(ctx, reqMsg)
	defer func() {
		logx.WithContext(ctx).WithDuration(time.Now().Sub(startTime)).
			Infof("PublishToDev msg:%v", reqMsg)
	}()
	err := n.client.Publish(ctx, fmt.Sprintf(topics.DeviceDownMsg, reqMsg.ProtocolCode, reqMsg.Handle, reqMsg.ProductID, reqMsg.DeviceName), devices.PublishToDev(
		reqMsg.Handle, reqMsg.Type, reqMsg.Payload, reqMsg.ProtocolCode,
		reqMsg.ProductID, reqMsg.DeviceName))

	if err != nil {
		logx.WithContext(ctx).Errorf("%s Publish failure err:%v", utils.FuncName(), err)
	}
	return err
}

func (n *pubDevClient) ReqToDeviceSync(ctx context.Context, reqMsg *deviceMsg.PublishMsg, timeout time.Duration, compareMsg CompareMsg) (
	payload []byte, err error) {
	start := time.Now()
	topic := fmt.Sprintf(topics.DeviceUpMsg, reqMsg.Handle, reqMsg.ProductID, reqMsg.DeviceName)
	done := make(chan struct{})
	sub, err := n.client.SubscribeWithID(topic, func(ctx context.Context, t time.Time, body []byte) error {
		msg, err := deviceMsg.GetDevPublish(ctx, body)
		if err != nil {
			logx.WithContext(ctx).Error(string(body), err)
			return err
		}
		if msg.Handle != reqMsg.Handle || msg.Type != reqMsg.Type { //不是订阅的topic
			return nil
		}
		if !compareMsg(msg.Payload) {
			return nil
		}
		now := time.Now().Sub(start)
		logx.WithContext(ctx).Infof("SubscribeWithID find use:%vms topic:%v msg:%v", now/time.Millisecond, topic, msg.String())
		payload = msg.Payload
		close(done)
		return nil
	})
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return nil, err
	}
	defer n.client.UnSubscribeWithID(topic, sub)
	err = n.PublishToDev(ctx, reqMsg)
	if err != nil {
		return nil, err
	}
	if timeout == 0 {
		timeout = time.Second * 5
	}
	select {
	case <-done:
		return
	case <-ctx.Done():
		return nil, errors.DeviceTimeOut
	case <-time.After(timeout):
		return nil, errors.DeviceTimeOut
	}
}
