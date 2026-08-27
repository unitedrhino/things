package devicemanagelogic

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNormalizeDeviceLifecycleStatus(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		row  relationDB.DmDeviceInfo
		want def.DeviceStatus
	}{
		{
			name: "expired owned device stays arrearage",
			row: relationDB.DmDeviceInfo{
				UserID:  2,
				Status:  def.DeviceStatusOnline,
				ExpTime: sql.NullTime{Time: now.Add(-time.Minute), Valid: true},
			},
			want: def.DeviceStatusArrearage,
		},
		{
			name: "unowned device is not forced to arrearage",
			row: relationDB.DmDeviceInfo{
				UserID:  def.RootNode,
				Status:  def.DeviceStatusOnline,
				ExpTime: sql.NullTime{Time: now.Add(-time.Minute), Valid: true},
			},
			want: def.DeviceStatusOnline,
		},
		{
			name: "future expiry preserves explicit state",
			row: relationDB.DmDeviceInfo{
				UserID:  2,
				Status:  def.DeviceStatusOffline,
				ExpTime: sql.NullTime{Time: now.Add(time.Minute), Valid: true},
			},
			want: def.DeviceStatusOffline,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			normalizeDeviceLifecycleStatus(&test.row, now)
			if test.row.Status != test.want {
				t.Fatalf("status = %d, want %d", test.row.Status, test.want)
			}
		})
	}
}

func TestAreaDeviceCountRefreshTargetsUseCapturedOldArea(t *testing.T) {
	newArea := &sys.AreaInfo{
		ProjectID:   102,
		AreaID:      12,
		AreaIDPath:  "12-",
		DeviceCount: wrapperspb.Int64(1),
		GroupCount:  wrapperspb.Int64(0),
	}

	targets := areaDeviceCountRefreshTargets(newArea, 101, 11, "11-")
	if len(targets) != 2 {
		t.Fatalf("area target count = %d, want 2", len(targets))
	}

	if targets[0] != newArea {
		t.Fatalf("first target should keep the fetched new area pointer")
	}

	oldArea := targets[1]
	if oldArea.ProjectID != 101 || oldArea.AreaID != 11 || oldArea.AreaIDPath != "11-" {
		t.Fatalf("old area target = project:%d area:%d path:%q, want project:101 area:11 path:%q",
			oldArea.ProjectID, oldArea.AreaID, oldArea.AreaIDPath, "11-")
	}
}

func TestProjectDeviceCountRefreshTargetsIncludeOldProject(t *testing.T) {
	targets := projectDeviceCountRefreshTargets(202, 201)
	want := []int64{202, 201}
	if !reflect.DeepEqual(targets, want) {
		t.Fatalf("project targets = %v, want %v", targets, want)
	}
}

func TestProjectDeviceCountRefreshTargetsSkipDuplicateAndInvalid(t *testing.T) {
	targets := projectDeviceCountRefreshTargets(201, 201)
	want := []int64{201}
	if !reflect.DeepEqual(targets, want) {
		t.Fatalf("duplicate project targets = %v, want %v", targets, want)
	}

	targets = projectDeviceCountRefreshTargets(0, 1)
	if len(targets) != 0 {
		t.Fatalf("invalid project targets = %v, want empty", targets)
	}
}
