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
1. 将RoleApp全局替换为模型的表名
2. 完善todo
*/

type RoleAppRepo struct {
	db *gorm.DB
}

func NewRoleAppRepo(in any) *RoleAppRepo {
	return &RoleAppRepo{db: stores.GetCommonConn(in)}
}

type RoleAppFilter struct {
	RoleID     int64
	TenantCode string
}

func (p RoleAppRepo) fmtFilter(ctx context.Context, f RoleAppFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if f.RoleID != 0 {
		db = db.Where("role_id =?", f.RoleID)
	}
	return db
}

func (p RoleAppRepo) Insert(ctx context.Context, data *SysRoleApp) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleAppRepo) FindOneByFilter(ctx context.Context, f RoleAppFilter) (*SysRoleApp, error) {
	var result SysRoleApp
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleAppRepo) FindByFilter(ctx context.Context, f RoleAppFilter, page *def.PageInfo) ([]*SysRoleApp, error) {
	var results []*SysRoleApp
	db := p.fmtFilter(ctx, f).Model(&SysRoleApp{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleAppRepo) CountByFilter(ctx context.Context, f RoleAppFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysRoleApp{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleAppRepo) Update(ctx context.Context, data *SysRoleApp) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleAppRepo) DeleteByFilter(ctx context.Context, f RoleAppFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysRoleApp{}).Error
	return stores.ErrFmt(err)
}

func (p RoleAppRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysRoleApp{}).Error
	return stores.ErrFmt(err)
}
func (p RoleAppRepo) FindOne(ctx context.Context, id int64) (*SysRoleApp, error) {
	var result SysRoleApp
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p RoleAppRepo) MultiInsert(ctx context.Context, data []*SysRoleApp) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysRoleApp{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p RoleAppRepo) MultiUpdate(ctx context.Context, roleID int64, appCodes []string) error {
	var datas []*SysRoleApp
	for _, v := range appCodes {
		datas = append(datas, &SysRoleApp{
			AppCode: v,
			RoleID:  roleID,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewRoleAppRepo(tx)
		err := rm.DeleteByFilter(ctx, RoleAppFilter{RoleID: roleID})
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
