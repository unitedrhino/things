package userdevicelogic

import (
	"context"
	"errors"
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/usershare"
	"gitee.com/unitedrhino/things/share/topics"
)

type recordingDeviceShareAcceptedPublisher struct {
	topic string
	arg   any
	err   error
}

func (p *recordingDeviceShareAcceptedPublisher) Publish(_ context.Context, topic string, arg any) error {
	p.topic = topic
	p.arg = arg
	return p.err
}

func TestShouldConsumeShareTokenAfterAccept(t *testing.T) {
	tests := []struct {
		name  string
		useBy string
		want  bool
	}{
		{name: "wechat single device token is one-time", useBy: "wechat_single_device", want: true},
		{name: "family token remains reusable", useBy: "family", want: false},
		{name: "empty useBy remains reusable", useBy: "", want: false},
		{name: "unknown useBy remains reusable", useBy: "batch_qr", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldConsumeShareTokenAfterAccept(tt.useBy)
			if got != tt.want {
				t.Fatalf("shouldConsumeShareTokenAfterAccept(%q) = %v, want %v", tt.useBy, got, tt.want)
			}
		})
	}
}

func TestBuildMultiShareTokenResponseIncludesLinkAndAuthExpiry(t *testing.T) {
	info := &dm.UserDeviceShareMultiInfo{
		CreatedTime: 1717500000,
		ExpTime:     1717600000,
	}

	got := buildMultiShareTokenResponse("share-token-1", info)

	if got.ShareToken != "share-token-1" {
		t.Fatalf("ShareToken = %q, want %q", got.ShareToken, "share-token-1")
	}
	wantLinkExpireAt := info.CreatedTime + int64(userShared.MultiDeviceShareTokenTTL.Seconds())
	if got.LinkExpireAt != wantLinkExpireAt {
		t.Fatalf("LinkExpireAt = %d, want %d", got.LinkExpireAt, wantLinkExpireAt)
	}
	if got.AuthExpireAt != info.ExpTime {
		t.Fatalf("AuthExpireAt = %d, want %d", got.AuthExpireAt, info.ExpTime)
	}
	if got.CreatedTime != info.CreatedTime {
		t.Fatalf("CreatedTime = %d, want %d", got.CreatedTime, info.CreatedTime)
	}
}

func TestBuildDeviceShareAcceptedEventUsesStableIDAndSortedDevices(t *testing.T) {
	in := &dm.UserDeviceShareMultiAcceptReq{
		ShareToken:        "share-token-1",
		SharedUserID:      202,
		SharedUserAccount: "receiver",
	}
	info := &dm.UserDeviceShareMultiInfo{
		UserID:    101,
		ProjectID: 303,
		UseBy:     "wechat_single_device",
	}
	devices := []*dm.DeviceShareInfo{
		{ProductID: "product-b", DeviceName: "device-2"},
		{ProductID: "product-a", DeviceName: "device-1"},
	}

	got := buildDeviceShareAcceptedEvent(in, info, "default", devices, 1717500000)
	reordered := buildDeviceShareAcceptedEvent(in, info, "default", []*dm.DeviceShareInfo{devices[1], devices[0]}, 1717500001)

	if got == nil {
		t.Fatal("buildDeviceShareAcceptedEvent() = nil, want event")
	}
	if got.EventID == "" {
		t.Fatal("EventID is empty")
	}
	if got.EventID != reordered.EventID {
		t.Fatalf("EventID changed after device order/timestamp changed: %q != %q", got.EventID, reordered.EventID)
	}
	if got.SharerUserID != info.UserID || got.ReceiverUserID != in.SharedUserID {
		t.Fatalf("unexpected users: sharer=%d receiver=%d", got.SharerUserID, got.ReceiverUserID)
	}
	if len(got.Devices) != 2 {
		t.Fatalf("len(Devices) = %d, want 2", len(got.Devices))
	}
	if got.Devices[0].ProductID != "product-a" || got.Devices[0].DeviceName != "device-1" {
		t.Fatalf("Devices[0] = %#v, want product-a/device-1", got.Devices[0])
	}
	if got.Devices[1].ProductID != "product-b" || got.Devices[1].DeviceName != "device-2" {
		t.Fatalf("Devices[1] = %#v, want product-b/device-2", got.Devices[1])
	}
}

func TestBuildDeviceShareAcceptedEventReturnsNilWithoutAcceptedDevices(t *testing.T) {
	got := buildDeviceShareAcceptedEvent(
		&dm.UserDeviceShareMultiAcceptReq{},
		&dm.UserDeviceShareMultiInfo{},
		"default",
		nil,
		1717500000,
	)

	if got != nil {
		t.Fatalf("buildDeviceShareAcceptedEvent() = %#v, want nil", got)
	}
}

func TestPublishDeviceShareAcceptedEventUsesDomainTopic(t *testing.T) {
	publisher := &recordingDeviceShareAcceptedPublisher{}
	event := &usershare.DeviceShareAcceptedEvent{EventID: "event-1"}

	err := publishDeviceShareAcceptedEvent(context.Background(), publisher, event)

	if err != nil {
		t.Fatalf("publishDeviceShareAcceptedEvent() error = %v", err)
	}
	if publisher.topic != topics.DmUserDeviceShareAccepted {
		t.Fatalf("topic = %q, want %q", publisher.topic, topics.DmUserDeviceShareAccepted)
	}
	if publisher.arg != event {
		t.Fatalf("arg = %#v, want event pointer", publisher.arg)
	}
}

func TestPublishDeviceShareAcceptedEventReturnsPublishError(t *testing.T) {
	wantErr := errors.New("publish failed")
	publisher := &recordingDeviceShareAcceptedPublisher{err: wantErr}

	err := publishDeviceShareAcceptedEvent(context.Background(), publisher, &usershare.DeviceShareAcceptedEvent{EventID: "event-1"})

	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
