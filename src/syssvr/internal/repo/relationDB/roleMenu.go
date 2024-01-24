package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleMenuRepo struct {
	db *gorm.DB
}

func NewRoleMenuRepo(in any) *RoleMenuRepo {
	return &RoleMenuRepo{db: stores.GetCommonConn(in)}
}

type RoleMenuFilter struct {
	TenantCode string
	RoleIDs    []int64
	AppCode    string
	ModuleCode string
}

func (p RoleMenuRepo) fmtFilter(ctx context.Context, f RoleMenuFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if f.ModuleCode != "" {
		db = db.Where("module_code =?", f.ModuleCode)
	}
	if len(f.RoleIDs) != 0 {
		db = db.Where("role_id in ?", f.RoleIDs)
	}
	if f.AppCode != "" {
		db = db.Where("app_code =?", f.AppCode)
	}
	return db
}

func (p RoleMenuRepo) Insert(ctx context.Context, data *SysRoleMenu) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleMenuRepo) FindOneByFilter(ctx context.Context, f RoleMenuFilter) (*SysRoleMenu, error) {
	var result SysRoleMenu
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleMenuRepo) FindByFilter(ctx context.Context, f RoleMenuFilter, page *def.PageInfo) ([]*SysRoleMenu, error) {
	var results []*SysRoleMenu
	db := p.fmtFilter(ctx, f).Model(&SysRoleMenu{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleMenuRepo) CountByFilter(ctx context.Context, f RoleMenuFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysRoleMenu{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleMenuRepo) Update(ctx context.Context, data *SysRoleMenu) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleMenuRepo) DeleteByFilter(ctx context.Context, f RoleMenuFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysRoleMenu{}).Error
	return stores.ErrFmt(err)
}

func (p RoleMenuRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysRoleMenu{}).Error
	return stores.ErrFmt(err)
}
func (p RoleMenuRepo) FindOne(ctx context.Context, id int64) (*SysRoleMenu, error) {
	var result SysRoleMenu
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p RoleMenuRepo) MultiInsert(ctx context.Context, data []*SysRoleMenu) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysRoleMenu{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p RoleMenuRepo) MultiUpdate(ctx context.Context, roleID int64, appCode string, moduleCode string, menuIDs []int64) error {
	var datas []*SysRoleMenu
	for _, v := range menuIDs {
		datas = append(datas, &SysRoleMenu{
			AppCode:    appCode,
			ModuleCode: moduleCode,
			RoleID:     roleID,
			MenuID:     v,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewRoleMenuRepo(tx)
		err := rm.DeleteByFilter(ctx, RoleMenuFilter{RoleIDs: []int64{roleID}, AppCode: appCode, ModuleCode: moduleCode})
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
