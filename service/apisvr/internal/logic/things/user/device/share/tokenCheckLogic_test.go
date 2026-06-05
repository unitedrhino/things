package share

import (
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

const testMultiDeviceShareTokenTTLSeconds int64 = 24 * 60 * 60

func TestToTokenCheckRespReturnsPublicTokenStateOnly(t *testing.T) {
	info := &dm.UserDeviceShareMultiInfo{
		Devices: []*dm.DeviceShareInfo{
			{ProductID: "p1", DeviceName: "d1"},
			{ProductID: "p2", DeviceName: "d2"},
		},
		CreatedTime: 1717500000,
		ExpTime:     1717600000,
		UseBy:       "wechat_single_device",
	}

	got := ToTokenCheckResp(info)

	if !got.Valid {
		t.Fatalf("Valid = false, want true")
	}
	if got.Reason != "" {
		t.Fatalf("Reason = %q, want empty", got.Reason)
	}
	wantLinkExpireAt := info.CreatedTime + testMultiDeviceShareTokenTTLSeconds
	if got.LinkExpireAt != wantLinkExpireAt {
		t.Fatalf("LinkExpireAt = %d, want %d", got.LinkExpireAt, wantLinkExpireAt)
	}
	if got.AuthExpireAt != info.ExpTime {
		t.Fatalf("AuthExpireAt = %d, want %d", got.AuthExpireAt, info.ExpTime)
	}
	if got.CreatedTime != info.CreatedTime {
		t.Fatalf("CreatedTime = %d, want %d", got.CreatedTime, info.CreatedTime)
	}
	if got.UseBy != info.UseBy {
		t.Fatalf("UseBy = %q, want %q", got.UseBy, info.UseBy)
	}
	if got.DeviceCount != 2 {
		t.Fatalf("DeviceCount = %d, want 2", got.DeviceCount)
	}
}

func TestInvalidTokenCheckRespUsesExpiredOrConsumedReason(t *testing.T) {
	got := InvalidTokenCheckResp()

	if got.Valid {
		t.Fatalf("Valid = true, want false")
	}
	if got.Reason != "expired_or_consumed" {
		t.Fatalf("Reason = %q, want expired_or_consumed", got.Reason)
	}
}
