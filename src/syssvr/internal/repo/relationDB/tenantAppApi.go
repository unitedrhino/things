package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantAppApiRepo struct {
	db *gorm.DB
}

func NewTenantAppApiRepo(in any) *TenantAppApiRepo {
	return &TenantAppApiRepo{db: stores.GetCommonConn(in)}
}

type TenantAppApiFilter struct {
	ApiIDs     []int64
	TenantCode string
	AppCode    string
	Route      string
	Method     string
	Group      string
	Name       string
	ModuleCode string
	WithRoles  bool
	IsNeedAuth int64
}

func (p TenantAppApiRepo) fmtFilter(ctx context.Context, f TenantAppApiFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithRoles {
		db = db.Preload("Roles")
	}
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if f.AppCode != "" {
		db = db.Where("app_code =?", f.AppCode)
	}
	if f.ApiIDs != nil {
		db = db.Where("id in ?", f.ApiIDs)
	}
	if f.IsNeedAuth != 0 {
		db = db.Where("is_need_auth =?", f.IsNeedAuth)
	}
	if f.Route != "" {
		db = db.Where("route like ?", "%"+f.Route+"%")
	}
	if f.ModuleCode != "" {
		db = db.Where("module_code =?", f.ModuleCode)
	}
	if f.Method != "" {
		db = db.Where("method = ?", f.Method)
	}
	if f.Group != "" {
		db = db.Where("group like ?", "%"+f.Group+"%")
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	return db
}

func (p TenantAppApiRepo) Insert(ctx context.Context, data *SysTenantAppApi) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantAppApiRepo) FindOneByFilter(ctx context.Context, f TenantAppApiFilter) (*SysTenantAppApi, error) {
	var result SysTenantAppApi
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p TenantAppApiRepo) FindByFilter(ctx context.Context, f TenantAppApiFilter, page *def.PageInfo) ([]*SysTenantAppApi, error) {
	var results []*SysTenantAppApi
	db := p.fmtFilter(ctx, f).Model(&SysTenantAppApi{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantAppApiRepo) CountByFilter(ctx context.Context, f TenantAppApiFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantAppApi{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantAppApiRepo) Update(ctx context.Context, data *SysTenantAppApi) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantAppApiRepo) DeleteByFilter(ctx context.Context, f TenantAppApiFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantAppApi{}).Error
	return stores.ErrFmt(err)
}
func (p TenantAppApiRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantAppApi{}).Error
	return stores.ErrFmt(err)
}

func (p TenantAppApiRepo) FindOne(ctx context.Context, id int64) (*SysTenantAppApi, error) {
	var result SysTenantAppApi
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantAppApiRepo) MultiInsert(ctx context.Context, data []*SysTenantAppApi) error {
	if len(data) == 0 {
		return nil
	}
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantAppApi{}).Create(data).Error
	return stores.ErrFmt(err)
}
