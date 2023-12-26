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
	RoleID int64
}

func (p RoleAppRepo) fmtFilter(ctx context.Context, f RoleAppFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.RoleID != 0 {
		db = db.Where("role_id =?", f.RoleID)
	}
	return db
}

func (p RoleAppRepo) Insert(ctx context.Context, data *SysTenantRoleApp) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleAppRepo) FindOneByFilter(ctx context.Context, f RoleAppFilter) (*SysTenantRoleApp, error) {
	var result SysTenantRoleApp
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleAppRepo) FindByFilter(ctx context.Context, f RoleAppFilter, page *def.PageInfo) ([]*SysTenantRoleApp, error) {
	var results []*SysTenantRoleApp
	db := p.fmtFilter(ctx, f).Model(&SysTenantRoleApp{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleAppRepo) CountByFilter(ctx context.Context, f RoleAppFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantRoleApp{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleAppRepo) Update(ctx context.Context, data *SysTenantRoleApp) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleAppRepo) DeleteByFilter(ctx context.Context, f RoleAppFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantRoleApp{}).Error
	return stores.ErrFmt(err)
}

func (p RoleAppRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantRoleApp{}).Error
	return stores.ErrFmt(err)
}
func (p RoleAppRepo) FindOne(ctx context.Context, id int64) (*SysTenantRoleApp, error) {
	var result SysTenantRoleApp
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p RoleAppRepo) MultiInsert(ctx context.Context, data []*SysTenantRoleApp) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantRoleApp{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p RoleAppRepo) MultiUpdate(ctx context.Context, roleID int64, appCodes []string) error {
	var datas []*SysTenantRoleApp
	for _, v := range appCodes {
		datas = append(datas, &SysTenantRoleApp{
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
