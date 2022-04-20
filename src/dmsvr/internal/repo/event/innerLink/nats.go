package innerLink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/ddsvr/ddExport"
	"github.com/i-Things/things/src/dmsvr/dmDef"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	deviceSend "github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	NatsClient struct {
		client *nats.Conn
	}
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
	return n.client.Publish(ddExport.TopicInnerPublish, ddExport.PublishToDev(ctx, topic, payload))
}

func (n *NatsClient) SubscribeDevSync(ctx context.Context, topic string) (*SubDev, error) {
	subscription, err := n.client.SubscribeSync(topic)
	if err != nil {
		return nil, err
	}
	return NewSubDev(subscription), nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	_, err := n.client.QueueSubscribe(ddExport.TopicDevPublishAll, dmDef.SvrName, func(msg *nats.Msg) {
		msg.Ack()
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ele, err := device.GetDevPublish(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Publish(ele)
		logx.WithContext(ctx).Info(ddExport.TopicDevPublishAll, msg.Subject, string(msg.Data), err)
	})
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(ddExport.TopicDevConnected, dmDef.SvrName, func(msg *nats.Msg) {
		msg.Ack()
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ele, err := device.GetDevConnMsg(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Connected(ele)
		logx.WithContext(ctx).Info(ddExport.TopicDevConnected, msg.Subject, string(msg.Data), err)
	})
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(ddExport.TopicDevDisconnected, dmDef.SvrName, func(msg *nats.Msg) {
		msg.Ack()
		emsg := events.GetEventMsg(msg.Data)
		if emsg == nil {
			logx.Error(msg.Subject, string(msg.Data))
			return
		}
		ctx := emsg.GetCtx()
		ele, err := device.GetDevConnMsg(ctx, emsg.GetData())
		if err != nil {
			logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
			return
		}
		err = handle(ctx).Disconnected(ele)
		logx.WithContext(ctx).Info(ddExport.TopicDevDisconnected, msg.Subject, string(msg.Data), err)
	})
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
	handle, err := n.SubscribeDevSync(ctx, fmt.Sprintf(ddExport.TopicDevPublish, productID, deviceName))
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
