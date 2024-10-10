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

type UserDeviceShareRepo struct {
	db *gorm.DB
}

func NewUserDeviceShareRepo(in any) *UserDeviceShareRepo {
	return &UserDeviceShareRepo{db: stores.GetCommonConn(in)}
}

type UserDeviceShareFilter struct {
	ProjectID     int64
	SharedUserIDs []int64
	Devices       []*devices.Core
	ProductID     string
	DeviceName    string
	SharedUserID  int64
	ID            int64
	IDs           []int64
	ExpTime       *stores.Cmp
}

func (p UserDeviceShareRepo) fmtFilter(ctx context.Context, f UserDeviceShareFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.ExpTime.Where(db, "exp_time")
	if f.SharedUserID != 0 {
		db = db.Where("shared_user_id = ?", f.SharedUserID)
	}
	if len(f.IDs) != 0 {
		db = db.Where("id in ?", f.IDs)
	}
	if f.ProjectID != 0 {
		db = db.Where("project_id=?", f.ProjectID)
	}
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if f.ProductID != "" {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name = ?", f.DeviceName)
	}
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
	return db
}

func (p UserDeviceShareRepo) Insert(ctx context.Context, data *DmUserDeviceShare) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserDeviceShareRepo) FindOneByFilter(ctx context.Context, f UserDeviceShareFilter) (*DmUserDeviceShare, error) {
	var result DmUserDeviceShare
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserDeviceShareRepo) FindByFilter(ctx context.Context, f UserDeviceShareFilter, page *stores.PageInfo) ([]*DmUserDeviceShare, error) {
	var results []*DmUserDeviceShare
	db := p.fmtFilter(ctx, f).Model(&DmUserDeviceShare{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserDeviceShareRepo) CountByFilter(ctx context.Context, f UserDeviceShareFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmUserDeviceShare{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserDeviceShareRepo) Update(ctx context.Context, data *DmUserDeviceShare) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d UserDeviceShareRepo) UpdateWithField(ctx context.Context, f UserDeviceShareFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmUserDeviceShare{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p UserDeviceShareRepo) DeleteByFilter(ctx context.Context, f UserDeviceShareFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmUserDeviceShare{}).Error
	return stores.ErrFmt(err)
}

func (p UserDeviceShareRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmUserDeviceShare{}).Error
	return stores.ErrFmt(err)
}
func (p UserDeviceShareRepo) FindOne(ctx context.Context, id int64) (*DmUserDeviceShare, error) {
	var result DmUserDeviceShare
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserDeviceShareRepo) MultiInsert(ctx context.Context, data []*DmUserDeviceShare) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmUserDeviceShare{}).Create(data).Error
	return stores.ErrFmt(err)
}
