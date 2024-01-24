package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleModuleRepo struct {
	db *gorm.DB
}

func NewRoleModuleRepo(in any) *RoleModuleRepo {
	return &RoleModuleRepo{db: stores.GetCommonConn(in)}
}

type RoleModuleFilter struct {
	TenantCode string
	RoleIDs    []int64
	AppCode    string
}

func (p RoleModuleRepo) fmtFilter(ctx context.Context, f RoleModuleFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if len(f.RoleIDs) != 0 {
		db = db.Where("role_id in ?", f.RoleIDs)
	}
	if f.AppCode != "" {
		db = db.Where("app_code =?", f.AppCode)
	}
	return db
}

func (p RoleModuleRepo) Insert(ctx context.Context, data *SysRoleModule) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleModuleRepo) FindOneByFilter(ctx context.Context, f RoleModuleFilter) (*SysRoleModule, error) {
	var result SysRoleModule
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleModuleRepo) FindByFilter(ctx context.Context, f RoleModuleFilter, page *def.PageInfo) ([]*SysRoleModule, error) {
	var results []*SysRoleModule
	db := p.fmtFilter(ctx, f).Model(&SysRoleModule{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleModuleRepo) CountByFilter(ctx context.Context, f RoleModuleFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysRoleModule{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleModuleRepo) Update(ctx context.Context, data *SysRoleModule) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleModuleRepo) DeleteByFilter(ctx context.Context, f RoleModuleFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysRoleModule{}).Error
	return stores.ErrFmt(err)
}

func (p RoleModuleRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysRoleModule{}).Error
	return stores.ErrFmt(err)
}
func (p RoleModuleRepo) FindOne(ctx context.Context, id int64) (*SysRoleModule, error) {
	var result SysRoleModule
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p RoleModuleRepo) MultiInsert(ctx context.Context, data []*SysRoleModule) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysRoleModule{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p RoleModuleRepo) MultiUpdate(ctx context.Context, roleID int64, appCode string, moduleCodes []string) error {
	var datas []*SysRoleModule
	for _, v := range moduleCodes {
		datas = append(datas, &SysRoleModule{
			AppCode:    appCode,
			RoleID:     roleID,
			ModuleCode: v,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewRoleModuleRepo(tx)
		err := rm.DeleteByFilter(ctx, RoleModuleFilter{RoleIDs: []int64{roleID}, AppCode: appCode})
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
