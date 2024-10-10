package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type DeviceTimingInfoRepo struct {
	db *gorm.DB
}

func NewDeviceTimerInfoRepo(in any) *DeviceTimingInfoRepo {
	return &DeviceTimingInfoRepo{db: stores.GetCommonConn(in)}
}

type DeviceTimerInfoFilter struct {
	Devices     []*devices.Core
	TriggerType string
	Status      int64
	ExecAt      *stores.Cmp
	LastRunTime *stores.Cmp
	Repeat      *stores.Cmp
}

func (p DeviceTimingInfoRepo) fmtFilter(ctx context.Context, f DeviceTimerInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.ExecAt.Where(db, "exec_at")
	db = f.LastRunTime.Where(db, "last_run_time")
	db = f.Repeat.Where(db, "exec_repeat")
	if len(f.Devices) != 0 {
		scope := func(db *gorm.DB) *gorm.DB {
			for i, d := range f.Devices {
				if i == 0 {
					db = db.Where("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
					continue
				}
				db = db.Or("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
			}
			return db
		}
		db = db.Where(scope(db))
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.TriggerType != "" {
		db = db.Where("trigger_type=?", f.TriggerType)
	}
	return db
}

func (p DeviceTimingInfoRepo) Insert(ctx context.Context, data *UdDeviceTimerInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeviceTimingInfoRepo) FindOneByFilter(ctx context.Context, f DeviceTimerInfoFilter) (*UdDeviceTimerInfo, error) {
	var result UdDeviceTimerInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DeviceTimingInfoRepo) FindByFilter(ctx context.Context, f DeviceTimerInfoFilter, page *stores.PageInfo) ([]*UdDeviceTimerInfo, error) {
	var results []*UdDeviceTimerInfo
	db := p.fmtFilter(ctx, f).Model(&UdDeviceTimerInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeviceTimingInfoRepo) CountByFilter(ctx context.Context, f DeviceTimerInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdDeviceTimerInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeviceTimingInfoRepo) Update(ctx context.Context, data *UdDeviceTimerInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeviceTimingInfoRepo) DeleteByFilter(ctx context.Context, f DeviceTimerInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdDeviceTimerInfo{}).Error
	return stores.ErrFmt(err)
}

func (p DeviceTimingInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdDeviceTimerInfo{}).Error
	return stores.ErrFmt(err)
}
func (p DeviceTimingInfoRepo) FindOne(ctx context.Context, id int64) (*UdDeviceTimerInfo, error) {
	var result UdDeviceTimerInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeviceTimingInfoRepo) MultiInsert(ctx context.Context, data []*UdDeviceTimerInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdDeviceTimerInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
