package logic

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type areaManageRecorder struct {
	areamanage.AreaManage
	updates []*sys.AreaInfo
}

func (m *areaManageRecorder) AreaInfoUpdate(ctx context.Context, in *areamanage.AreaInfo, opts ...grpc.CallOption) (*areamanage.Empty, error) {
	m.updates = append(m.updates, &sys.AreaInfo{
		ProjectID:   in.ProjectID,
		AreaID:      in.AreaID,
		DeviceCount: in.DeviceCount,
	})
	return &sys.Empty{}, nil
}

func TestDirectFillAreaDeviceCountRefreshesAncestorsWithLeafCountCollision(t *testing.T) {
	ctx := ctxs.WithRoot(ctxs.WithProjectID(context.Background(), 100))
	initAreaCountTestDB(t, ctx)

	insertDeviceForAreaCount(t, ctx, "active-leaf", 11, "10-11-", 0)
	insertDeviceForAreaCount(t, ctx, "active-other-child", 12, "10-12-", 0)
	insertDeviceForAreaCount(t, ctx, "deleted-leaf", 11, "10-11-", stores.DeletedTime(time.Now().Unix()))

	areaM := &areaManageRecorder{}
	svcCtx := &svc.ServiceContext{AreaM: areaM}
	leafArea := &sys.AreaInfo{
		ProjectID:   100,
		AreaID:      11,
		AreaIDPath:  "10-11-",
		DeviceCount: wrapperspb.Int64(2),
	}

	if err := DirectFillAreaDeviceCount(ctx, svcCtx, 0, leafArea); err != nil {
		t.Fatalf("DirectFillAreaDeviceCount error: %v", err)
	}

	got := map[int64]int64{}
	for _, update := range areaM.updates {
		got[update.AreaID] = update.DeviceCount.GetValue()
	}

	if got[10] != 2 {
		t.Fatalf("parent area update count = %d, want 2; all updates: %v", got[10], got)
	}
	if got[11] != 1 {
		t.Fatalf("leaf area update count = %d, want 1 active device; all updates: %v", got[11], got)
	}
}

func initAreaCountTestDB(t *testing.T, ctx context.Context) {
	t.Helper()

	stores.InitConn(conf.Database{
		DBType: conf.Sqlite,
		DSN:    filepath.Join(t.TempDir(), "area-count.db"),
	})
	db := stores.GetCommonConn(ctx)
	if err := db.Migrator().DropTable(&relationDB.DmDeviceInfo{}); err != nil {
		t.Fatalf("drop dm_device_info: %v", err)
	}
	if err := db.AutoMigrate(&relationDB.DmDeviceInfo{}); err != nil {
		t.Fatalf("migrate dm_device_info: %v", err)
	}
}

func insertDeviceForAreaCount(t *testing.T, ctx context.Context, deviceName string, areaID int64, areaIDPath string, deletedTime stores.DeletedTime) {
	t.Helper()

	err := stores.GetCommonConn(ctx).Exec(`
insert into dm_device_info (
	tenant_code, project_id, area_id, area_id_path, product_id, device_name,
	device_alias, secret, cert, imei, mac, version, hard_info, soft_info,
	tags, schema_alias, protocol_conf, sub_protocol_conf, deleted_time
) values (?, ?, ?, ?, ?, ?, '', '', '', '', '', '', '', '', '{}', '{}', '{}', '{}', ?)`,
		"default", 100, areaID, areaIDPath, "test-product", deviceName, deletedTime,
	).Error
	if err != nil {
		t.Fatalf("insert device %s: %v", deviceName, err)
	}
}
