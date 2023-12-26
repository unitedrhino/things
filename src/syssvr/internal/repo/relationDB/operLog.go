package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type OperLogRepo struct {
	db *gorm.DB
}

func NewOperLogRepo(in any) *OperLogRepo {
	return &OperLogRepo{db: stores.GetCommonConn(in)}
}

type OperLogFilter struct {
	OperName     string
	OperUserName string
	BusinessType int64
}

func (p OperLogRepo) fmtFilter(ctx context.Context, f OperLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.OperName != "" {
		db = db.Where("`oper_name` like ?", "%"+f.OperName+"%")
	}
	if f.OperUserName != "" {
		db = db.Where("`oper_user_name` like ?", "%"+f.OperUserName+"%")
	}
	if f.BusinessType > 0 {
		db = db.Where("`business_type`= ?", f.BusinessType)
	}
	return db
}

func (p OperLogRepo) Insert(ctx context.Context, data *SysTenantOperLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OperLogRepo) FindOneByFilter(ctx context.Context, f OperLogFilter) (*SysTenantOperLog, error) {
	var result SysTenantOperLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OperLogRepo) FindByFilter(ctx context.Context, f OperLogFilter, page *def.PageInfo) ([]*SysTenantOperLog, error) {
	var results []*SysTenantOperLog
	db := p.fmtFilter(ctx, f).Model(&SysTenantOperLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OperLogRepo) CountByFilter(ctx context.Context, f OperLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantOperLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OperLogRepo) Update(ctx context.Context, data *SysTenantOperLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OperLogRepo) DeleteByFilter(ctx context.Context, f OperLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantOperLog{}).Error
	return stores.ErrFmt(err)
}

func (p OperLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantOperLog{}).Error
	return stores.ErrFmt(err)
}
func (p OperLogRepo) FindOne(ctx context.Context, id int64) (*SysTenantOperLog, error) {
	var result SysTenantOperLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
