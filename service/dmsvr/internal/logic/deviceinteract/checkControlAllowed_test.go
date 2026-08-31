package deviceinteractlogic

import (
	"testing"
	"time"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCheckDeviceControlAllowed(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	tests := []struct {
		name    string
		device  *dm.DeviceInfo
		wantErr bool
	}{
		{
			name: "expired owned device",
			device: &dm.DeviceInfo{
				UserID:  2,
				Status:  def.DeviceStatusOnline,
				ExpTime: wrapperspb.Int64(now.Add(-time.Second).Unix()),
			},
			wantErr: true,
		},
		{
			name: "expiry boundary",
			device: &dm.DeviceInfo{
				UserID:  2,
				Status:  def.DeviceStatusOnline,
				ExpTime: wrapperspb.Int64(now.Unix()),
			},
			wantErr: true,
		},
		{
			name: "future expiry",
			device: &dm.DeviceInfo{
				UserID:  2,
				Status:  def.DeviceStatusArrearage,
				ExpTime: wrapperspb.Int64(now.Add(time.Second).Unix()),
			},
		},
		{
			name: "no expiry",
			device: &dm.DeviceInfo{
				UserID: 2,
				Status: def.DeviceStatusArrearage,
			},
		},
		{
			name: "unowned historical expiry",
			device: &dm.DeviceInfo{
				UserID:  def.RootNode,
				Status:  def.DeviceStatusOnline,
				ExpTime: wrapperspb.Int64(now.Add(-time.Second).Unix()),
			},
		},
		{
			name: "nil device",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkDeviceControlAllowed(test.device, now)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected expired device control error")
				}
				if got := errors.Fmt(err).GetCode(); got != errors.NotEnable.GetCode() {
					t.Fatalf("error code = %d, want %d", got, errors.NotEnable.GetCode())
				}
				if got := errors.Fmt(err).GetMsg(); got != expiredDeviceControlMessage {
					t.Fatalf("error message = %q, want %q", got, expiredDeviceControlMessage)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
