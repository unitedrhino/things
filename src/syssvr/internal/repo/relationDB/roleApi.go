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

func NewRoleApiRepo(in any) *RoleApiRepo {
	return &RoleApiRepo{db: stores.GetCommonConn(in)}
}

type RoleApiFilter struct {
	RoleIDs []int64
	AppCode string
}

func (p RoleApiRepo) fmtFilter(ctx context.Context, f RoleApiFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.RoleIDs) != 0 {
		db = db.Where("role_id in ?", f.RoleIDs)
	}
	if f.AppCode != "" {
		db = db.Where("app_code =?", f.AppCode)
	}
	return db
}

func (p RoleApiRepo) Insert(ctx context.Context, data *SysTenantRoleApi) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p RoleApiRepo) FindOneByFilter(ctx context.Context, f RoleApiFilter) (*SysTenantRoleApi, error) {
	var result SysTenantRoleApi
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleApiRepo) FindByFilter(ctx context.Context, f RoleApiFilter, page *def.PageInfo) ([]*SysTenantRoleApi, error) {
	var results []*SysTenantRoleApi
	db := p.fmtFilter(ctx, f).Model(&SysTenantRoleApi{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p RoleApiRepo) CountByFilter(ctx context.Context, f RoleApiFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantRoleApi{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p RoleApiRepo) Update(ctx context.Context, data *SysTenantRoleApi) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p RoleApiRepo) DeleteByFilter(ctx context.Context, f RoleApiFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantRoleApi{}).Error
	return stores.ErrFmt(err)
}

func (p RoleApiRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantRoleApi{}).Error
	return stores.ErrFmt(err)
}
func (p RoleApiRepo) FindOne(ctx context.Context, id int64) (*SysTenantRoleApi, error) {
	var result SysTenantRoleApi
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p RoleApiRepo) MultiInsert(ctx context.Context, data []*SysTenantRoleApi) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantRoleApi{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p RoleApiRepo) MultiUpdate(ctx context.Context, roleID int64, appCode string, apiIDs []int64) error {
	var datas []*SysTenantRoleApi
	for _, v := range apiIDs {
		datas = append(datas, &SysTenantRoleApi{
			AppCode: appCode,
			RoleID:  roleID,
			ApiID:   v,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewRoleApiRepo(tx)
		err := rm.DeleteByFilter(ctx, RoleApiFilter{RoleIDs: []int64{roleID}, AppCode: appCode})
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
