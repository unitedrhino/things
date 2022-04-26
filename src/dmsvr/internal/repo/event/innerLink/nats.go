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
	deviceSend "github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type (
	NatsClient struct {
		client nats.JetStreamContext
	}
)

//topic 定义
const (
	ThingsStreamName = "thing_msg"
	// TopicDevPublish dd模块收到设备的发布消息后向内部推送以下topic 最后两个是产品id和设备名称
	TopicDevPublish    = "dd.thing.device.clients.publish.%s.%s"
	TopicDevPublishAll = "dd.thing.device.clients.publish.>"

	// TopicDevConnected dd模块收到设备的登录消息后向内部推送以下topic
	TopicDevConnected = "dd.thing.device.clients.connected"
	// TopicDevDisconnected dd模块收到设备的登出消息后向内部推送以下topic
	TopicDevDisconnected = "dd.thing.device.clients.disconnected"
	// TopicInnerPublish dd模块订阅以下topic,收到内部的发布消息后向设备推送
	TopicInnerPublish = "dd.thing.inner.publish"
	TopicThing        = "dd.thing.device.clients.>"
)
const (
	ThingsDeliverGroup     = "things_dm_group"
	ThingsQueueConsumeName = "things_dm_queue_consume"
	ThingsAllConsumeName   = "things_dm_all_consume"
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
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name: ThingsStreamName,
		Subjects: []string{
			TopicThing,
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = js.AddConsumer(ThingsStreamName, &nats.ConsumerConfig{
		Durable:   ThingsQueueConsumeName,
		AckPolicy: nats.AckExplicitPolicy,
		//MaxRequestBatch:   10,
		//MaxRequestExpires: 2 * time.Second,
		DeliverPolicy:  nats.DeliverLastPolicy,
		DeliverSubject: nats.NewInbox(),
		DeliverGroup:   ThingsDeliverGroup,
	})
	_, err = js.AddConsumer(ThingsStreamName, &nats.ConsumerConfig{
		Durable:   ThingsQueueConsumeName + "2",
		AckPolicy: nats.AckExplicitPolicy,
		//MaxRequestBatch:   10,
		//MaxRequestExpires: 2 * time.Second,
		DeliverPolicy:  nats.DeliverLastPolicy,
		DeliverSubject: nats.NewInbox(),
		DeliverGroup:   ThingsDeliverGroup + "2",
	})
	_, err = js.AddConsumer(ThingsStreamName, &nats.ConsumerConfig{
		Durable:   ThingsAllConsumeName,
		AckPolicy: nats.AckExplicitPolicy,
		//MaxRequestBatch:   10,
		//MaxRequestExpires: 2 * time.Second,
		DeliverPolicy:  nats.DeliverNewPolicy,
		DeliverSubject: nats.NewInbox(),
	})
	if err != nil {
		return nil, err
	}
	return &NatsClient{client: js}, nil
}

func (n *NatsClient) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	_, err := n.client.Publish(TopicInnerPublish, devices.PublishToDev(ctx, topic, payload))
	return err
}

func (n *NatsClient) SubscribeDevSync(ctx context.Context, topic string) (*SubDev, error) {
	subscription, err := n.client.SubscribeSync(topic, nats.Durable(ThingsAllConsumeName),
		nats.BindStream(ThingsStreamName))
	if err != nil {
		return nil, err
	}
	return NewSubDev(subscription), nil
}

func (n *NatsClient) Subscribe(handle Handle) error {
	//_, err := n.client.QueueSubscribe(TopicDevPublishAll, ThingsDeliverGroup, func(msg *nats.Msg) {
	//	msg.Ack()
	//	emsg := events.GetEventMsg(msg.Data)
	//	if emsg == nil {
	//		logx.Error(msg.Subject, string(msg.Data))
	//		return
	//	}
	//	ctx := emsg.GetCtx()
	//	ele, err := device.GetDevPublish(ctx, emsg.GetData())
	//	if err != nil {
	//		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
	//		return
	//	}
	//	err = handle(ctx).Publish(ele)
	//	logx.WithContext(ctx).Infof("%s|topic:%v,subject:%v,data:%v,err:%v", utils.FuncName(),
	//		TopicDevPublishAll, msg.Subject, string(msg.Data), err)
	//}, nats.Durable(ThingsQueueConsumeName), nats.BindStream(ThingsStreamName))
	//if err != nil {
	//	return err
	//}
	//_, err = n.client.QueueSubscribe(TopicDevPublishAll, ThingsDeliverGroup+"2", func(msg *nats.Msg) {
	//	err := msg.Ack()
	//	if err != nil {
	//		logx.Error(msg.Subject, string(msg.Data), err)
	//		return
	//	}
	//	emsg := events.GetEventMsg(msg.Data)
	//	if emsg == nil {
	//		logx.Error(msg.Subject, string(msg.Data))
	//		return
	//	}
	//	ctx := emsg.GetCtx()
	//	ele, err := device.GetDevPublish(ctx, emsg.GetData())
	//	if err != nil {
	//		logx.WithContext(ctx).Error(msg.Subject, string(msg.Data), err)
	//		return
	//	}
	//	err = handle(ctx).Publish(ele)
	//	logx.WithContext(ctx).Infof("%s|topic:%v,subject:%v,data:%v,err:%v", utils.FuncName(),
	//		TopicDevPublishAll, msg.Subject, string(msg.Data), err)
	//}, nats.Durable(ThingsQueueConsumeName+"2"), nats.BindStream(ThingsStreamName))
	//if err != nil {
	//	return err
	//}
	_, err := n.client.QueueSubscribe(TopicDevConnected, ThingsDeliverGroup, func(msg *nats.Msg) {
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
		logx.WithContext(ctx).Infof("%s|topic:%v,subject:%v,data:%v,err:%v", utils.FuncName(),
			TopicDevConnected, msg.Subject, string(msg.Data), err)
	}, nats.Durable(ThingsQueueConsumeName), nats.BindStream(ThingsStreamName))
	if err != nil {
		return err
	}
	_, err = n.client.QueueSubscribe(TopicDevDisconnected, ThingsDeliverGroup, func(msg *nats.Msg) {
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
		logx.WithContext(ctx).Infof("%s|topic:%v,subject:%v,data:%v,err:%v", utils.FuncName(),
			TopicDevDisconnected, msg.Subject, string(msg.Data), err)
	}, nats.Durable(ThingsQueueConsumeName), nats.BindStream(ThingsStreamName))
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
	handle, err := n.SubscribeDevSync(ctx, fmt.Sprintf(TopicDevPublish, productID, deviceName))
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
