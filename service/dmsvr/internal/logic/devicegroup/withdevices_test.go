package devicegrouplogic

import (
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func TestFillGroupDevices_FillsMatchingGroupsRecursively(t *testing.T) {
	groups := []*dm.GroupInfo{
		{
			Id: 1,
			Children: []*dm.GroupInfo{
				{Id: 2},
			},
		},
		{Id: 3},
	}
	groupDevices := []*relationDB.DmGroupDevice{
		{
			GroupID: 1,
			Device: &relationDB.DmDeviceInfo{ProductID: "p1", DeviceName: "d1"},
		},
		{
			GroupID: 2,
			Device: &relationDB.DmDeviceInfo{ProductID: "p2", DeviceName: "d2"},
		},
		{
			GroupID: 3,
			Device: &relationDB.DmDeviceInfo{ProductID: "p3", DeviceName: "d3"},
		},
		{
			GroupID: 99,
			Device: &relationDB.DmDeviceInfo{ProductID: "px", DeviceName: "dx"},
		},
		{
			GroupID: 1,
		},
	}

	fillGroupDevices(groups, groupDevices)

	assertDeviceCore(t, groups[0].Devices, "p1", "d1")
	assertDeviceCore(t, groups[0].Children[0].Devices, "p2", "d2")
	assertDeviceCore(t, groups[1].Devices, "p3", "d3")
}

func assertDeviceCore(t *testing.T, devices []*dm.DeviceCore, productID, deviceName string) {
	t.Helper()
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	if devices[0].ProductID != productID || devices[0].DeviceName != deviceName {
		t.Fatalf("unexpected device core: %+v", devices[0])
	}
}
