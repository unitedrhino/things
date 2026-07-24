package userdevicelogic

import (
	"testing"

	shareerrors "gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

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

func TestValidateMultiAcceptDevicesRejectsUnboundDevice(t *testing.T) {
	tokenDevices := []*dm.DeviceShareInfo{
		{ProductID: "product-1", DeviceName: "device-1"},
	}
	requestedDevices := []*dm.DeviceCore{
		{ProductID: "product-1", DeviceName: "device-1"},
	}

	err := validateMultiAcceptDevices(101, tokenDevices, requestedDevices, func(_ string, _ string) (int64, error) {
		return 2, nil
	})

	if !shareerrors.Cmp(err, shareerrors.DeviceNotBound) {
		t.Fatalf("error = %v, want DeviceNotBound", err)
	}
}

func TestValidateMultiAcceptDevicesRejectsDeviceMovedToAnotherProject(t *testing.T) {
	tokenDevices := []*dm.DeviceShareInfo{
		{ProductID: "product-1", DeviceName: "device-1"},
	}
	requestedDevices := []*dm.DeviceCore{
		{ProductID: "product-1", DeviceName: "device-1"},
	}

	err := validateMultiAcceptDevices(101, tokenDevices, requestedDevices, func(_ string, _ string) (int64, error) {
		return 202, nil
	})

	if !shareerrors.Cmp(err, shareerrors.DeviceNotBound) {
		t.Fatalf("error = %v, want DeviceNotBound", err)
	}
}

func TestValidateMultiAcceptDevicesAllowsDeviceInOriginalProject(t *testing.T) {
	tokenDevices := []*dm.DeviceShareInfo{
		{ProductID: "product-1", DeviceName: "device-1"},
	}
	requestedDevices := []*dm.DeviceCore{
		{ProductID: "product-1", DeviceName: "device-1"},
	}

	err := validateMultiAcceptDevices(101, tokenDevices, requestedDevices, func(_ string, _ string) (int64, error) {
		return 101, nil
	})

	if err != nil {
		t.Fatalf("validateMultiAcceptDevices() error = %v, want nil", err)
	}
}

func TestPrepareMultiShareInfoStoresCreatorAndProject(t *testing.T) {
	info := &dm.UserDeviceShareMultiInfo{}

	prepareMultiShareInfo(info, 101, 202, 1717500000)

	if info.UserID != 101 {
		t.Fatalf("UserID = %d, want 101", info.UserID)
	}
	if info.ProjectID != 202 {
		t.Fatalf("ProjectID = %d, want 202", info.ProjectID)
	}
	if info.CreatedTime != 1717500000 {
		t.Fatalf("CreatedTime = %d, want 1717500000", info.CreatedTime)
	}
}
