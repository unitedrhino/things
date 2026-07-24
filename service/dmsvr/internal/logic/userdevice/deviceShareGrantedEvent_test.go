// 本文件用于独立验证分享授权生效事件的统一契约和发布主题。
package userdevicelogic

import (
	"context"
	"errors"
	"testing"

	"gitee.com/unitedrhino/things/share/domain/usershare"
	"gitee.com/unitedrhino/things/share/topics"
)

// focusedDeviceShareGrantedPublisher 记录测试期间发布的事件。
type focusedDeviceShareGrantedPublisher struct {
	topic string
	arg   any
	err   error
}

// Publish 记录主题和消息，并返回预设错误。
func (p *focusedDeviceShareGrantedPublisher) Publish(_ context.Context, topic string, arg any) error {
	p.topic = topic
	p.arg = arg
	return p.err
}

// TestDeviceShareGrantedEventStableIDAndSources 验证两类分享来源和幂等编号。
func TestDeviceShareGrantedEventStableIDAndSources(t *testing.T) {
	devices := []usershare.DeviceShareGrantedDevice{
		{ShareID: 402, ProductID: "product-b", DeviceName: "device-2"},
		{ShareID: 401, ProductID: "product-a", DeviceName: "device-1"},
	}
	wechat := buildDeviceShareGrantedEvent(
		usershare.DeviceShareGrantSourceWechatAccept,
		"share-token-1",
		101,
		202,
		"receiver",
		303,
		"default",
		"wechat_single_device",
		devices,
		1717500000,
	)
	reordered := buildDeviceShareGrantedEvent(
		usershare.DeviceShareGrantSourceWechatAccept,
		"share-token-1",
		101,
		202,
		"receiver",
		303,
		"default",
		"wechat_single_device",
		[]usershare.DeviceShareGrantedDevice{devices[1], devices[0]},
		1717500001,
	)
	account := buildDeviceShareGrantedEvent(
		usershare.DeviceShareGrantSourceAccountDirect,
		"",
		101,
		202,
		"receiver",
		303,
		"default",
		"",
		devices[:1],
		1717500000,
	)

	if wechat == nil || reordered == nil || account == nil {
		t.Fatal("buildDeviceShareGrantedEvent() returned nil")
	}
	if wechat.EventID != reordered.EventID {
		t.Fatalf("EventID changed after device order/timestamp changed: %q != %q", wechat.EventID, reordered.EventID)
	}
	if wechat.Source != usershare.DeviceShareGrantSourceWechatAccept {
		t.Fatalf("wechat Source = %q", wechat.Source)
	}
	if account.Source != usershare.DeviceShareGrantSourceAccountDirect || account.ShareToken != "" {
		t.Fatalf("account event = %#v", account)
	}
	if wechat.Devices[0].ShareID != 401 || wechat.Devices[1].ShareID != 402 {
		t.Fatalf("Devices = %#v, want stable order", wechat.Devices)
	}
}

// TestDeviceShareGrantedEventReturnsNilWithoutDevices 验证空授权记录不会生成事件。
func TestDeviceShareGrantedEventReturnsNilWithoutDevices(t *testing.T) {
	got := buildDeviceShareGrantedEvent("", "", 0, 0, "", 0, "", "", nil, 1717500000)
	if got != nil {
		t.Fatalf("buildDeviceShareGrantedEvent() = %#v, want nil", got)
	}
}

// TestDeviceShareGrantedEventPublisher 验证统一主题及发布错误透传。
func TestDeviceShareGrantedEventPublisher(t *testing.T) {
	event := &usershare.DeviceShareGrantedEvent{EventID: "event-1"}
	publisher := &focusedDeviceShareGrantedPublisher{}

	err := publishDeviceShareGrantedEvent(context.Background(), publisher, event)
	if err != nil {
		t.Fatalf("publishDeviceShareGrantedEvent() error = %v", err)
	}
	if publisher.topic != topics.DmUserDeviceShareGranted || publisher.arg != event {
		t.Fatalf("published topic=%q arg=%#v", publisher.topic, publisher.arg)
	}

	wantErr := errors.New("publish failed")
	publisher.err = wantErr
	err = publishDeviceShareGrantedEvent(context.Background(), publisher, event)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
