package pubApp

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/topics"
)

type (
	PubApp interface {
		DeviceStatusConnected(ctx context.Context, di *dm.DeviceInfo, msg application.ConnectMsg) error
		DeviceStatusDisConnected(ctx context.Context, di *dm.DeviceInfo, msg application.ConnectMsg) error
		//应用事件通知-设备物模型属性上报通知
		DeviceThingPropertyReport(ctx context.Context, di *dm.DeviceInfo, msg application.PropertyReport) error
		DeviceThingPropertyReportV2(ctx context.Context, di *dm.DeviceInfo, msg application.PropertyReportV2) error
		DeviceThingActionReport(ctx context.Context, di *dm.DeviceInfo, msg application.ActionReport) error
		DeviceThingEventReport(ctx context.Context, di *dm.DeviceInfo, msg application.EventReport) error
	}
	pubApp struct {
		client *eventBus.FastEvent
	}
)

func NewPubApp(f *eventBus.FastEvent) (PubApp, error) {
	return &pubApp{client: f}, nil
}

func (n *pubApp) DeviceThingPropertyReportV2(ctx context.Context, di *dm.DeviceInfo, msg application.PropertyReportV2) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingPropertyV2, di.TenantCode, msg.Device.ProductID, msg.Device.DeviceName)
	err := n.client.Publish(ctx, topic, msg)
	return err
}

// 应用事件通知-设备物模型属性上报通知
func (n *pubApp) DeviceThingPropertyReport(ctx context.Context, di *dm.DeviceInfo, msg application.PropertyReport) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingProperty, di.TenantCode, msg.Device.ProductID, msg.Device.DeviceName, msg.Identifier)
	err := n.client.Publish(ctx, topic, msg)
	return err
}

func (n *pubApp) DeviceThingActionReport(ctx context.Context, di *dm.DeviceInfo, msg application.ActionReport) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingAction, di.TenantCode,
		msg.Device.ProductID, msg.Device.DeviceName, msg.ActionID, msg.ReqType, msg.Dir)
	err := n.client.Publish(ctx, topic, msg)
	return err
}

func (n *pubApp) DeviceThingEventReport(ctx context.Context, di *dm.DeviceInfo, msg application.EventReport) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingEvent, di.TenantCode,
		msg.Device.ProductID, msg.Device.DeviceName, msg.Type, msg.Identifier)
	err := n.client.Publish(ctx, topic, msg)
	return err
}

func (n *pubApp) DeviceStatusConnected(ctx context.Context, di *dm.DeviceInfo, msg application.ConnectMsg) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusConnected, di.TenantCode, msg.Device.ProductID, msg.Device.DeviceName)
	err := n.client.Publish(ctx, topic, msg)
	return err
}

func (n *pubApp) DeviceStatusDisConnected(ctx context.Context, di *dm.DeviceInfo, msg application.ConnectMsg) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusDisConnected, di.TenantCode, msg.Device.ProductID, msg.Device.DeviceName)
	err := n.client.Publish(ctx, topic, msg)
	return err
}
