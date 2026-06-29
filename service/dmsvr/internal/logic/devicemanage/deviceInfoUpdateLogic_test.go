package devicemanagelogic

import (
	"reflect"
	"testing"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
