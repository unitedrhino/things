package relationDB

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"github.com/glebarez/sqlite"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func newDeviceStatusTestRepo(t *testing.T) DeviceInfoRepo {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Exec(`CREATE TABLE dm_device_info (
		id INTEGER PRIMARY KEY,
		tenant_code TEXT,
		project_id INTEGER,
		area_id INTEGER,
		area_id_path TEXT,
		product_id TEXT,
		device_name TEXT,
		status INTEGER NOT NULL,
		is_online INTEGER NOT NULL,
		exp_time DATETIME,
		user_id INTEGER NOT NULL,
		updated_by INTEGER NOT NULL DEFAULT 0,
		updated_time DATETIME,
		distributor_updated_time DATETIME,
		deleted_time INTEGER NOT NULL DEFAULT 0
	)`).Error; err != nil {
		t.Fatal(err)
	}
	return DeviceInfoRepo{db: db}
}

func deviceStatusTestContext() context.Context {
	return ctxs.SetUserCtx(context.Background(), &ctxs.UserCtx{
		TenantCode: "default",
		ProjectID:  100,
		IsAdmin:    true,
		InnerCtx: ctxs.InnerCtx{
			AllTenant:  true,
			AllProject: true,
			AllArea:    true,
		},
	})
}

func insertDeviceStatusTestRow(t *testing.T, repo DeviceInfoRepo, id int64, status def.DeviceStatus, online int64, expTime any, userID int64) {
	t.Helper()
	if err := repo.db.Exec(`INSERT INTO dm_device_info
		(id, tenant_code, project_id, area_id, area_id_path, product_id, device_name, status, is_online, exp_time, user_id, deleted_time)
		VALUES (?, 'default', 1, 1, '1', 'p', ?, ?, ?, ?, ?, 0)`, id, fmt.Sprintf("d%d", id), status, online, expTime, userID).Error; err != nil {
		t.Fatal(err)
	}
}

func TestUpdateExpiredDeviceStatusBatchesAndNoop(t *testing.T) {
	repo := newDeviceStatusTestRepo(t)
	ctx := deviceStatusTestContext()
	cutoff := time.Now()
	for id := int64(1); id <= 250; id++ {
		insertDeviceStatusTestRow(t, repo, id, def.DeviceStatusOnline, def.True, cutoff.Add(-time.Hour), 2)
	}
	insertDeviceStatusTestRow(t, repo, 251, def.DeviceStatusArrearage, def.True, cutoff.Add(-time.Hour), 2)
	insertDeviceStatusTestRow(t, repo, 252, def.DeviceStatusOnline, def.True, cutoff.Add(time.Hour), 2)
	insertDeviceStatusTestRow(t, repo, 253, def.DeviceStatusOnline, def.True, cutoff.Add(-time.Hour), 1)

	updated, batches, err := repo.UpdateExpiredDeviceStatus(ctx, cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if updated != 250 || batches != 2 {
		t.Fatalf("updated=%d batches=%d, want 250/2", updated, batches)
	}

	updated, batches, err = repo.UpdateExpiredDeviceStatus(ctx, cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if updated != 0 || batches != 0 {
		t.Fatalf("second run updated=%d batches=%d, want zero write", updated, batches)
	}
}

func TestDeviceStatusBatchRechecksCurrentStatus(t *testing.T) {
	repo := newDeviceStatusTestRepo(t)
	ctx := deviceStatusTestContext()
	insertDeviceStatusTestRow(t, repo, 1, def.DeviceStatusAbnormal, def.True, nil, 2)
	insertDeviceStatusTestRow(t, repo, 2, def.DeviceStatusOnline, def.True, nil, 2)
	insertDeviceStatusTestRow(t, repo, 3, def.DeviceStatusOffline, def.False, nil, 2)

	recovered, err := repo.RecoverAbnormalDeviceStatusBatch(ctx, []int64{2, 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(recovered) != 1 || recovered[0].ID != 1 {
		t.Fatalf("recovered ids = %v, want [1]", deviceStatusIDs(recovered))
	}
	if err = repo.db.Exec("UPDATE dm_device_info SET status = ? WHERE id = ?", def.DeviceStatusArrearage, 1).Error; err != nil {
		t.Fatal(err)
	}

	marked, err := repo.MarkAbnormalDeviceStatusBatch(ctx, []int64{3, 1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if got := deviceStatusIDs(marked); fmt.Sprint(got) != "[2 3]" {
		t.Fatalf("marked ids = %v, want [2 3]", got)
	}
}

func TestDeviceStatusBatchLimit(t *testing.T) {
	repo := newDeviceStatusTestRepo(t)
	ids := make([]int64, DeviceStatusUpdateBatchSize+1)
	_, err := repo.MarkAbnormalDeviceStatusBatch(deviceStatusTestContext(), ids)
	if err == nil {
		t.Fatal("expected batch limit error")
	}
}

func TestRetryDeviceStatusUpdateOnlyForDeadlocks(t *testing.T) {
	attempts := 0
	err := retryDeviceStatusUpdate(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return &mysql.MySQLError{Number: 1213, Message: "deadlock"}
		}
		return nil
	})
	if err != nil || attempts != 3 {
		t.Fatalf("deadlock retry err=%v attempts=%d", err, attempts)
	}

	attempts = 0
	errExpected := errors.New("ordinary error")
	err = retryDeviceStatusUpdate(context.Background(), func() error {
		attempts++
		return errExpected
	})
	if !errors.Is(err, errExpected) || attempts != 1 {
		t.Fatalf("ordinary error retry err=%v attempts=%d", err, attempts)
	}
}

func deviceStatusIDs(devices []*DmDeviceInfo) []int64 {
	ids := make([]int64, 0, len(devices))
	for _, device := range devices {
		ids = append(ids, device.ID)
	}
	return ids
}
