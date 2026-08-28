package relationDB

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/def"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDeviceStatusUpdatesConcurrentOnMySQL(t *testing.T) {
	dsn := os.Getenv("DM_DEVICE_INFO_MYSQL_DSN")
	if dsn == "" {
		t.Skip("DM_DEVICE_INFO_MYSQL_DSN is not set")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	var database string
	if err = db.Raw("SELECT DATABASE()").Scan(&database).Error; err != nil {
		t.Fatal(err)
	}
	if database != "things_p0_deadlock_test" {
		t.Fatalf("refusing destructive integration test in database %q", database)
	}
	t.Cleanup(func() { _ = db.Exec("DROP TABLE IF EXISTS dm_device_info").Error })
	if err = db.Exec(`CREATE TABLE dm_device_info (
		id BIGINT PRIMARY KEY,
		tenant_code VARCHAR(50) NOT NULL,
		project_id BIGINT NOT NULL,
		area_id BIGINT NOT NULL,
		area_id_path VARCHAR(100) NOT NULL,
		product_id VARCHAR(100) NOT NULL,
		device_name VARCHAR(100) NOT NULL,
		status SMALLINT NOT NULL,
		is_online SMALLINT NOT NULL,
		exp_time DATETIME NULL,
		user_id BIGINT NOT NULL,
		updated_by BIGINT NOT NULL DEFAULT 0,
		updated_time DATETIME NULL,
		distributor_updated_time DATETIME NULL,
		deleted_time BIGINT NOT NULL DEFAULT 0,
		KEY idx_exp_time (exp_time, user_id, deleted_time)
	) ENGINE=InnoDB`).Error; err != nil {
		t.Fatal(err)
	}

	repo := DeviceInfoRepo{db: db}
	ctx := deviceStatusTestContext()
	cutoff := time.Now()
	ids := make([]int64, DeviceStatusUpdateBatchSize)
	for i := range ids {
		ids[i] = int64(i + 1)
	}

	for round := 0; round < 20; round++ {
		if err = db.Exec("DELETE FROM dm_device_info").Error; err != nil {
			t.Fatal(err)
		}
		for _, id := range ids {
			if err = db.Exec(`INSERT INTO dm_device_info
				(id, tenant_code, project_id, area_id, area_id_path, product_id, device_name, status, is_online, exp_time, user_id, deleted_time)
				VALUES (?, 'default', 1, 1, '1', 'p', ?, ?, ?, ?, 2, 0)`,
				id, fmt.Sprintf("d%d", id), def.DeviceStatusAbnormal, def.True, cutoff.Add(-time.Hour)).Error; err != nil {
				t.Fatal(err)
			}
		}

		start := make(chan struct{})
		errCh := make(chan error, 2)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			_, _, updateErr := repo.UpdateExpiredDeviceStatus(ctx, cutoff)
			errCh <- updateErr
		}()
		go func() {
			defer wg.Done()
			<-start
			_, updateErr := repo.RecoverAbnormalDeviceStatusBatch(ctx, ids)
			errCh <- updateErr
		}()
		close(start)
		wg.Wait()
		close(errCh)
		for updateErr := range errCh {
			if updateErr != nil {
				t.Fatalf("round %d concurrent update failed: %v", round, updateErr)
			}
		}
	}
}
