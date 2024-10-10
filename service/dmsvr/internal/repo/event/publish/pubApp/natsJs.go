package pubApp

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/domain/application"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/share/events/topics"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	NatsJsClient struct {
		client nats.JetStreamContext
	}
)

func newNatsJsClient(conf conf.NatsConf) (*NatsJsClient, error) {
	nc, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		logx.Errorf("Nats 连接失败 err:%v", err)
		return nil, err
	}
	return &NatsJsClient{client: nc}, nil
}

// 应用事件通知-设备物模型属性上报通知
func (n *NatsJsClient) DeviceThingPropertyReport(ctx context.Context, msg application.PropertyReport) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingProperty, msg.Device.ProductID, msg.Device.DeviceName, msg.Identifier)
	_, err := n.client.Publish(topic, pubMsg)
	return err
}

func (n *NatsJsClient) DeviceThingActionReport(ctx context.Context, msg application.ActionReport) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingAction,
		msg.Device.ProductID, msg.Device.DeviceName, msg.ActionID, msg.ReqType, msg.Dir)
	_, err := n.client.Publish(topic, pubMsg)
	return err
}

func (n *NatsJsClient) DeviceThingEventReport(ctx context.Context, msg application.EventReport) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingEvent,
		msg.Device.ProductID, msg.Device.DeviceName, msg.Type, msg.Identifier)
	_, err := n.client.Publish(topic, pubMsg)
	return err
}

func (n *NatsJsClient) DeviceStatusConnected(ctx context.Context, msg application.ConnectMsg) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusConnected, msg.Device.ProductID, msg.Device.DeviceName)
	_, err := n.client.Publish(topic, pubMsg)
	return err
}

func (n *NatsJsClient) DeviceStatusDisConnected(ctx context.Context, msg application.ConnectMsg) error {
	data, _ := json.Marshal(msg)
	pubMsg := events.NewEventMsg(ctx, data)
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusDisConnected, msg.Device.ProductID, msg.Device.DeviceName)
	_, err := n.client.Publish(topic, pubMsg)
	return err
}
