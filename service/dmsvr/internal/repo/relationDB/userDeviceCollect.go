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

type UserCollectDeviceRepo struct {
	db *gorm.DB
}

func NewUserDeviceCollectRepo(in any) *UserCollectDeviceRepo {
	return &UserCollectDeviceRepo{db: stores.GetCommonConn(in)}
}

type UserDeviceCollectFilter struct {
	Cores  []*devices.Core
	UserID int64
}

func (p UserCollectDeviceRepo) fmtFilter(ctx context.Context, f UserDeviceCollectFilter) *gorm.DB {
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

func (p UserCollectDeviceRepo) Insert(ctx context.Context, data *DmUserCollectDevice) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserCollectDeviceRepo) FindOneByFilter(ctx context.Context, f UserDeviceCollectFilter) (*DmUserCollectDevice, error) {
	var result DmUserCollectDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserCollectDeviceRepo) FindByFilter(ctx context.Context, f UserDeviceCollectFilter, page *def.PageInfo) ([]*DmUserCollectDevice, error) {
	var results []*DmUserCollectDevice
	db := p.fmtFilter(ctx, f).Model(&DmUserCollectDevice{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserCollectDeviceRepo) CountByFilter(ctx context.Context, f UserDeviceCollectFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmUserCollectDevice{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserCollectDeviceRepo) Update(ctx context.Context, data *DmUserCollectDevice) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserCollectDeviceRepo) DeleteByFilter(ctx context.Context, f UserDeviceCollectFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmUserCollectDevice{}).Error
	return stores.ErrFmt(err)
}

func (p UserCollectDeviceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmUserCollectDevice{}).Error
	return stores.ErrFmt(err)
}
func (p UserCollectDeviceRepo) FindOne(ctx context.Context, id int64) (*DmUserCollectDevice, error) {
	var result DmUserCollectDevice
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserCollectDeviceRepo) MultiInsert(ctx context.Context, data []*DmUserCollectDevice) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmUserCollectDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}
