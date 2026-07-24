package cache

import (
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func TestMultiShareContainsDevice(t *testing.T) {
	info := &dm.UserDeviceShareMultiInfo{
		Devices: []*dm.DeviceShareInfo{
			{ProductID: "product-1", DeviceName: "device-1"},
			{ProductID: "product-2", DeviceName: "device-2"},
		},
	}

	if !multiShareContainsDevice(info, "product-2", "device-2") {
		t.Fatal("multiShareContainsDevice() = false, want true")
	}
	if multiShareContainsDevice(info, "product-2", "device-1") {
		t.Fatal("multiShareContainsDevice() = true for mismatched device, want false")
	}
}
