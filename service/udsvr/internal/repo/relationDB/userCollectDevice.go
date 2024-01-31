package relationDB

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type UserCollectDeviceRepo struct {
	db *gorm.DB
}

func NewUserCollectDeviceRepo(in any) *UserCollectDeviceRepo {
	return &UserCollectDeviceRepo{db: stores.GetCommonConn(in)}
}

type UserCollectDeviceFilter struct {
	Cores  []*devices.Core
	UserID int64
}

func (p UserCollectDeviceRepo) fmtFilter(ctx context.Context, f UserCollectDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	if len(f.Cores) != 0 {
		scope := func(db *gorm.DB) *gorm.DB {
			for i, d := range f.Cores {
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
	return db
}

func (p UserCollectDeviceRepo) Insert(ctx context.Context, data *UdUserCollectDevice) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserCollectDeviceRepo) FindOneByFilter(ctx context.Context, f UserCollectDeviceFilter) (*UdUserCollectDevice, error) {
	var result UdUserCollectDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserCollectDeviceRepo) FindByFilter(ctx context.Context, f UserCollectDeviceFilter, page *def.PageInfo) ([]*UdUserCollectDevice, error) {
	var results []*UdUserCollectDevice
	db := p.fmtFilter(ctx, f).Model(&UdUserCollectDevice{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserCollectDeviceRepo) CountByFilter(ctx context.Context, f UserCollectDeviceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdUserCollectDevice{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserCollectDeviceRepo) Update(ctx context.Context, data *UdUserCollectDevice) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserCollectDeviceRepo) DeleteByFilter(ctx context.Context, f UserCollectDeviceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdUserCollectDevice{}).Error
	return stores.ErrFmt(err)
}

func (p UserCollectDeviceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdUserCollectDevice{}).Error
	return stores.ErrFmt(err)
}
func (p UserCollectDeviceRepo) FindOne(ctx context.Context, id int64) (*UdUserCollectDevice, error) {
	var result UdUserCollectDevice
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserCollectDeviceRepo) MultiInsert(ctx context.Context, data []*UdUserCollectDevice) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdUserCollectDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}
