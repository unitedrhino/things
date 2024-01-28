package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantApiRepo struct {
	db *gorm.DB
}

func NewTenantAccessRepo(in any) *TenantApiRepo {
	return &TenantApiRepo{db: stores.GetCommonConn(in)}
}

type TenantAccessFilter struct {
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

func (p TenantApiRepo) fmtFilter(ctx context.Context, f TenantAccessFilter) *gorm.DB {
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

func (p TenantApiRepo) Insert(ctx context.Context, data *SysTenantAccess) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantApiRepo) FindOneByFilter(ctx context.Context, f TenantAccessFilter) (*SysTenantAccess, error) {
	var result SysTenantAccess
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p TenantApiRepo) FindByFilter(ctx context.Context, f TenantAccessFilter, page *def.PageInfo) ([]*SysTenantAccess, error) {
	var results []*SysTenantAccess
	db := p.fmtFilter(ctx, f).Model(&SysTenantAccess{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantApiRepo) CountByFilter(ctx context.Context, f TenantAccessFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantAccess{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantApiRepo) Update(ctx context.Context, data *SysTenantAccess) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantApiRepo) DeleteByFilter(ctx context.Context, f TenantAccessFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantAccess{}).Error
	return stores.ErrFmt(err)
}
func (p TenantApiRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantAccess{}).Error
	return stores.ErrFmt(err)
}

func (p TenantApiRepo) FindOne(ctx context.Context, id int64) (*SysTenantAccess, error) {
	var result SysTenantAccess
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantApiRepo) MultiInsert(ctx context.Context, data []*SysTenantAccess) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantAccess{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p TenantApiRepo) MultiUpdate(ctx context.Context, tenantCode string, AccessCodes []string) error {
	var datas []*SysTenantAccess
	for _, v := range AccessCodes {
		datas = append(datas, &SysTenantAccess{
			TenantCode: stores.TenantCode(tenantCode),
			AccessCode: v,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewTenantAccessRepo(tx)
		err := rm.DeleteByFilter(ctx, TenantAccessFilter{TenantCode: tenantCode})
		if err != nil {
			return err
		}
		if len(datas) != 0 {
			err = rm.MultiInsert(ctx, datas)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}
