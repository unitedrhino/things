package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
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

type UserDeviceShareRepo struct {
	db *gorm.DB
}

func NewUserDeviceShareRepo(in any) *UserDeviceShareRepo {
	return &UserDeviceShareRepo{db: stores.GetCommonConn(in)}
}

type UserDeviceShareFilter struct {
	ProductID  string
	DeviceName string
	UserID     int64
	ID         int64
}

func (p UserDeviceShareRepo) fmtFilter(ctx context.Context, f UserDeviceShareFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
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
func (p UserDeviceShareRepo) FindByFilter(ctx context.Context, f UserDeviceShareFilter, page *def.PageInfo) ([]*DmUserDeviceShare, error) {
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
