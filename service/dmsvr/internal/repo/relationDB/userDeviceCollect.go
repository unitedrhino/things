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

type UserDeviceCollectRepo struct {
	db *gorm.DB
}

func NewUserDeviceCollectRepo(in any) *UserDeviceCollectRepo {
	return &UserDeviceCollectRepo{db: stores.GetCommonConn(in)}
}

type UserDeviceCollectFilter struct {
	Cores  []*devices.Core
	UserID int64
}

func (p UserDeviceCollectRepo) fmtFilter(ctx context.Context, f UserDeviceCollectFilter) *gorm.DB {
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

func (p UserDeviceCollectRepo) Insert(ctx context.Context, data *DmUserDeviceCollect) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserDeviceCollectRepo) FindOneByFilter(ctx context.Context, f UserDeviceCollectFilter) (*DmUserDeviceCollect, error) {
	var result DmUserDeviceCollect
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserDeviceCollectRepo) FindByFilter(ctx context.Context, f UserDeviceCollectFilter, page *stores.PageInfo) ([]*DmUserDeviceCollect, error) {
	var results []*DmUserDeviceCollect
	db := p.fmtFilter(ctx, f).Model(&DmUserDeviceCollect{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserDeviceCollectRepo) CountByFilter(ctx context.Context, f UserDeviceCollectFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmUserDeviceCollect{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserDeviceCollectRepo) Update(ctx context.Context, data *DmUserDeviceCollect) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d UserDeviceCollectRepo) UpdateWithField(ctx context.Context, f UserDeviceCollectFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmUserDeviceCollect{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p UserDeviceCollectRepo) DeleteByFilter(ctx context.Context, f UserDeviceCollectFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmUserDeviceCollect{}).Error
	return stores.ErrFmt(err)
}

func (p UserDeviceCollectRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmUserDeviceCollect{}).Error
	return stores.ErrFmt(err)
}
func (p UserDeviceCollectRepo) FindOne(ctx context.Context, id int64) (*DmUserDeviceCollect, error) {
	var result DmUserDeviceCollect
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserDeviceCollectRepo) MultiInsert(ctx context.Context, data []*DmUserDeviceCollect) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmUserDeviceCollect{}).Create(data).Error
	return stores.ErrFmt(err)
}
