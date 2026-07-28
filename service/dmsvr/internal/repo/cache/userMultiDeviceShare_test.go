package cache

import (
	"reflect"
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

func TestMultiShareIndexKeysIncludeCreatorAndUniqueDevices(t *testing.T) {
	manager := &UserMultiDeviceShareManager{}
	info := &dm.UserDeviceShareMultiInfo{
		UserID: 101,
		Devices: []*dm.DeviceShareInfo{
			{ProductID: "product-1", DeviceName: "device-1"},
			{ProductID: "product-1", DeviceName: "device-1"},
			{ProductID: "product-2", DeviceName: "device-2"},
		},
	}

	got := manager.multiShareIndexKeys("tenant-a", info)
	want := []string{
		"things:device:share:batch:list:tenant-a:101",
		"things:device:share:batch:device:tenant-a:product-1:device-1",
		"things:device:share:batch:device:tenant-a:product-2:device-2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("multiShareIndexKeys() = %#v, want %#v", got, want)
	}
}
