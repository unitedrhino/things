package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantAppMenuRepo struct {
	db *gorm.DB
}

func NewTenantAppMenuRepo(in any) *TenantAppMenuRepo {
	return &TenantAppMenuRepo{db: stores.GetCommonConn(in)}
}

type TenantAppMenuFilter struct {
	ModuleCode string
	Name       string
	TenantCode string
	AppCode    string
	Path       string
	MenuIDs    []int64
}

func (p TenantAppMenuRepo) fmtFilter(ctx context.Context, f TenantAppMenuFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if f.AppCode != "" {
		db = db.Where("app_code =?", f.AppCode)
	}
	if f.ModuleCode != "" {
		db = db.Where("module_code =?", f.ModuleCode)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Path != "" {
		db = db.Where("path like ?", "%"+f.Path+"%")
	}
	if len(f.MenuIDs) != 0 {
		db = db.Where("id in ?", f.MenuIDs)
	}
	return db
}

func (p TenantAppMenuRepo) Insert(ctx context.Context, data *SysTenantAppMenu) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantAppMenuRepo) FindOneByFilter(ctx context.Context, f TenantAppMenuFilter) (*SysTenantAppMenu, error) {
	var result SysTenantAppMenu
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantAppMenuRepo) FindByFilter(ctx context.Context, f TenantAppMenuFilter, page *def.PageInfo) ([]*SysTenantAppMenu, error) {
	var results []*SysTenantAppMenu
	db := p.fmtFilter(ctx, f).Model(&SysTenantAppMenu{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantAppMenuRepo) CountByFilter(ctx context.Context, f TenantAppMenuFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantAppMenu{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantAppMenuRepo) Update(ctx context.Context, data *SysTenantAppMenu) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantAppMenuRepo) DeleteByFilter(ctx context.Context, f TenantAppMenuFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantAppMenu{}).Error
	return stores.ErrFmt(err)
}

func (p TenantAppMenuRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantAppMenu{}).Error
	return stores.ErrFmt(err)
}
func (p TenantAppMenuRepo) FindOne(ctx context.Context, id int64) (*SysTenantAppMenu, error) {
	var result SysTenantAppMenu
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantAppMenuRepo) MultiInsert(ctx context.Context, data []*SysTenantAppMenu) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantAppMenu{}).Create(data).Error
	return stores.ErrFmt(err)
}
