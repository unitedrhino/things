package devicemanagelogic

import (
	"context"
	"testing"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
)

type recordedTransferShareUpdate struct {
	filter relationDB.UserDeviceShareFilter
	fields map[string]any
}

type recordingTransferShareRepo struct {
	shares        []*relationDB.DmUserDeviceShare
	deleteFilters []relationDB.UserDeviceShareFilter
	updates       []recordedTransferShareUpdate
}

func newRecordingTransferShareRepo() *recordingTransferShareRepo {
	return &recordingTransferShareRepo{
		shares: []*relationDB.DmUserDeviceShare{
			{
				ProductID:    "product-1",
				DeviceName:   "device-1",
				SharedUserID: 202,
			},
			{
				ProductID:    "product-1",
				DeviceName:   "device-1",
				SharedUserID: 404,
			},
		},
	}
}

func (r *recordingTransferShareRepo) FindByFilter(
	_ context.Context,
	_ relationDB.UserDeviceShareFilter,
	_ *stores.PageInfo,
) ([]*relationDB.DmUserDeviceShare, error) {
	return r.shares, nil
}

func (r *recordingTransferShareRepo) DeleteByFilter(
	_ context.Context,
	filter relationDB.UserDeviceShareFilter,
) error {
	r.deleteFilters = append(r.deleteFilters, filter)
	return nil
}

func (r *recordingTransferShareRepo) UpdateWithField(
	_ context.Context,
	filter relationDB.UserDeviceShareFilter,
	fields map[string]any,
) error {
	r.updates = append(r.updates, recordedTransferShareUpdate{
		filter: filter,
		fields: fields,
	})
	return nil
}

func testTransferDevices() []*devices.Core {
	return []*devices.Core{
		{ProductID: "product-1", DeviceName: "device-1"},
	}
}

func TestShouldPreserveDeviceTransferSharesOnlyForNoCleanData(t *testing.T) {
	tests := []struct {
		name        string
		isCleanData int64
		want        bool
	}{
		{name: "clean data removes shares", isCleanData: def.True, want: false},
		{name: "no clean data preserves shares", isCleanData: def.False, want: true},
		{name: "missing value keeps safe legacy behavior", isCleanData: 0, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := shouldPreserveDeviceTransferShares(test.isCleanData); got != test.want {
				t.Fatalf("shouldPreserveDeviceTransferShares(%d) = %v, want %v", test.isCleanData, got, test.want)
			}
		})
	}
}

func TestApplyDeviceTransferSharePolicyDeletesAllWhenCleaning(t *testing.T) {
	repo := newRecordingTransferShareRepo()
	shares, err := applyDeviceTransferSharePolicy(
		context.Background(),
		repo,
		testTransferDevices(),
		202,
		303,
		"tenant-b",
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 2 {
		t.Fatalf("len(shares) = %d, want 2", len(shares))
	}
	if len(repo.deleteFilters) != 1 || repo.deleteFilters[0].SharedUserID != 0 {
		t.Fatalf("deleteFilters = %#v, want delete all device shares", repo.deleteFilters)
	}
	if len(repo.updates) != 0 {
		t.Fatalf("updates = %#v, want none", repo.updates)
	}
}

func TestApplyDeviceTransferSharePolicyKeepsOthersAndRemovesRecipient(t *testing.T) {
	repo := newRecordingTransferShareRepo()
	shares, err := applyDeviceTransferSharePolicy(
		context.Background(),
		repo,
		testTransferDevices(),
		202,
		303,
		"tenant-b",
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 2 {
		t.Fatalf("len(shares) = %d, want 2", len(shares))
	}
	if len(repo.deleteFilters) != 1 || repo.deleteFilters[0].SharedUserID != 202 {
		t.Fatalf("deleteFilters = %#v, want recipient-only delete", repo.deleteFilters)
	}
	if len(repo.updates) != 1 {
		t.Fatalf("updates = %#v, want one preserved-share update", repo.updates)
	}
	if repo.updates[0].filter.ExcludeUserID != 202 {
		t.Fatalf("ExcludeUserID = %d, want 202", repo.updates[0].filter.ExcludeUserID)
	}
	if repo.updates[0].fields["project_id"] != int64(303) {
		t.Fatalf("project_id = %#v, want 303", repo.updates[0].fields["project_id"])
	}
	if repo.updates[0].fields["tenant_code"] != "tenant-b" {
		t.Fatalf("tenant_code = %#v, want tenant-b", repo.updates[0].fields["tenant_code"])
	}
}
