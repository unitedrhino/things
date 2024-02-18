package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
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

func NewDeviceTimingInfoRepo(in any) *DeviceTimingInfoRepo {
	return &DeviceTimingInfoRepo{db: stores.GetCommonConn(in)}
}

type DeviceTimingInfoFilter struct {
	Devices     []*devices.Core
	TriggerType string
	Status      int64
}

func (p DeviceTimingInfoRepo) fmtFilter(ctx context.Context, f DeviceTimingInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
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

func (p DeviceTimingInfoRepo) Insert(ctx context.Context, data *UdDeviceTimingInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeviceTimingInfoRepo) FindOneByFilter(ctx context.Context, f DeviceTimingInfoFilter) (*UdDeviceTimingInfo, error) {
	var result UdDeviceTimingInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DeviceTimingInfoRepo) FindByFilter(ctx context.Context, f DeviceTimingInfoFilter, page *def.PageInfo) ([]*UdDeviceTimingInfo, error) {
	var results []*UdDeviceTimingInfo
	db := p.fmtFilter(ctx, f).Model(&UdDeviceTimingInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeviceTimingInfoRepo) CountByFilter(ctx context.Context, f DeviceTimingInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdDeviceTimingInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeviceTimingInfoRepo) Update(ctx context.Context, data *UdDeviceTimingInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeviceTimingInfoRepo) DeleteByFilter(ctx context.Context, f DeviceTimingInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdDeviceTimingInfo{}).Error
	return stores.ErrFmt(err)
}

func (p DeviceTimingInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdDeviceTimingInfo{}).Error
	return stores.ErrFmt(err)
}
func (p DeviceTimingInfoRepo) FindOne(ctx context.Context, id int64) (*UdDeviceTimingInfo, error) {
	var result UdDeviceTimingInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeviceTimingInfoRepo) MultiInsert(ctx context.Context, data []*UdDeviceTimingInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdDeviceTimingInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
