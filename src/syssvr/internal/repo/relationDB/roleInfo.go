package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type RoleInfoRepo struct {
	db *gorm.DB
}

func NewRoleInfoRepo(in any) *RoleInfoRepo {
	return &RoleInfoRepo{db: stores.GetCommonConn(in)}
}

type RoleInfoFilter struct {
	IDs         []int64
	WithAppInfo bool
	AppCode     string
	Name        string
	Status      int64
	TenantCode  string
}

func (p RoleInfoRepo) fmtFilter(ctx context.Context, f RoleInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code=?", f.TenantCode)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.AppCode != "" {
		db = db.Where("app_code=?", f.AppCode)
	}
	if len(f.IDs) > 0 {
		db = db.Where("id in ?", f.IDs)
	}
	if f.WithAppInfo {
		db = db.Preload("Apps")
	}
	if f.Status > 0 {
		db = db.Where("status= ?", f.Status)
	}
	return db
}

func (p RoleInfoRepo) Insert(ctx context.Context, data *SysTenantRoleInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleInfoRepo) FindOneByFilter(ctx context.Context, f RoleInfoFilter) (*SysTenantRoleInfo, error) {
	var result SysTenantRoleInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleInfoRepo) FindByFilter(ctx context.Context, f RoleInfoFilter, page *def.PageInfo) ([]*SysTenantRoleInfo, error) {
	var results []*SysTenantRoleInfo
	db := p.fmtFilter(ctx, f).Model(&SysTenantRoleInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleInfoRepo) CountByFilter(ctx context.Context, f RoleInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantRoleInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleInfoRepo) Update(ctx context.Context, data *SysTenantRoleInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleInfoRepo) DeleteByFilter(ctx context.Context, f RoleInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantRoleInfo{}).Error
	return stores.ErrFmt(err)
}

func (p RoleInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantRoleInfo{}).Error
	return stores.ErrFmt(err)
}
func (p RoleInfoRepo) FindOne(ctx context.Context, id int64) (*SysTenantRoleInfo, error) {
	var result SysTenantRoleInfo
	db := p.db.WithContext(ctx)
	err := db.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
