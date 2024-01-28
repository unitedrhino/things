package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将RoleApi全局替换为模型的表名
2. 完善todo
*/

type RoleApiRepo struct {
	db *gorm.DB
}

func NewRoleAccessRepo(in any) *RoleApiRepo {
	return &RoleApiRepo{db: stores.GetCommonConn(in)}
}

type RoleAccessFilter struct {
	TenantCode  string
	RoleIDs     []int64
	AccessCodes []string
}

func (p RoleApiRepo) fmtFilter(ctx context.Context, f RoleAccessFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if len(f.AccessCodes) != 0 {
		db = db.Where("api_scope_code =?", f.AccessCodes)
	}
	if len(f.RoleIDs) != 0 {
		db = db.Where("role_id in ?", f.RoleIDs)
	}
	return db
}

func (p RoleApiRepo) Insert(ctx context.Context, data *SysRoleAccess) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleApiRepo) FindOneByFilter(ctx context.Context, f RoleAccessFilter) (*SysRoleAccess, error) {
	var result SysRoleAccess
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleApiRepo) FindByFilter(ctx context.Context, f RoleAccessFilter, page *def.PageInfo) ([]*SysRoleAccess, error) {
	var results []*SysRoleAccess
	db := p.fmtFilter(ctx, f).Model(&SysRoleAccess{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleApiRepo) CountByFilter(ctx context.Context, f RoleAccessFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysRoleAccess{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleApiRepo) Update(ctx context.Context, data *SysRoleAccess) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleApiRepo) DeleteByFilter(ctx context.Context, f RoleAccessFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysRoleAccess{}).Error
	return stores.ErrFmt(err)
}

func (p RoleApiRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysRoleAccess{}).Error
	return stores.ErrFmt(err)
}
func (p RoleApiRepo) FindOne(ctx context.Context, id int64) (*SysRoleAccess, error) {
	var result SysRoleAccess
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p RoleApiRepo) MultiInsert(ctx context.Context, data []*SysRoleAccess) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysRoleAccess{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p RoleApiRepo) MultiUpdate(ctx context.Context, roleID int64, AccessCodes []string) error {
	var datas []*SysRoleAccess
	for _, v := range AccessCodes {
		datas = append(datas, &SysRoleAccess{
			RoleID:     roleID,
			AccessCode: v,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewRoleAccessRepo(tx)
		err := rm.DeleteByFilter(ctx, RoleAccessFilter{RoleIDs: []int64{roleID}})
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
