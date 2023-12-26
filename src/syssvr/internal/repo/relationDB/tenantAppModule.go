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
1. 将TenantAppModule全局替换为模型的表名
2. 完善todo
*/

type TenantAppModuleRepo struct {
	db *gorm.DB
}

func NewTenantAppModuleRepo(in any) *TenantAppModuleRepo {
	return &TenantAppModuleRepo{db: stores.GetCommonConn(in)}
}

type TenantAppModuleFilter struct {
	ID          int64
	AppCodes    []string
	AppCode     string
	ModuleCodes []string
	TenantCode  string
}

func (p TenantAppModuleRepo) fmtFilter(ctx context.Context, f TenantAppModuleFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.AppCode != "" {
		db = db.Where("app_code like ?", "%"+f.AppCode+"%")
	}
	if f.ID != 0 {
		db = db.Where("id =?", f.ID)
	}
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if len(f.AppCodes) != 0 {
		db = db.Where("app_code in ?", f.AppCodes)
	}
	return db
}

func (p TenantAppModuleRepo) Insert(ctx context.Context, data *SysTenantAppModule) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantAppModuleRepo) FindOneByFilter(ctx context.Context, f TenantAppModuleFilter) (*SysTenantAppModule, error) {
	var result SysTenantAppModule
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantAppModuleRepo) FindByFilter(ctx context.Context, f TenantAppModuleFilter, page *def.PageInfo) ([]*SysTenantAppModule, error) {
	var results []*SysTenantAppModule
	db := p.fmtFilter(ctx, f).Model(&SysTenantAppModule{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantAppModuleRepo) CountByFilter(ctx context.Context, f TenantAppModuleFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantAppModule{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantAppModuleRepo) Update(ctx context.Context, data *SysTenantAppModule) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantAppModuleRepo) DeleteByFilter(ctx context.Context, f TenantAppModuleFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantAppModule{}).Error
	return stores.ErrFmt(err)
}

func (p TenantAppModuleRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantAppModule{}).Error
	return stores.ErrFmt(err)
}
func (p TenantAppModuleRepo) FindOne(ctx context.Context, id int64) (*SysTenantAppModule, error) {
	var result SysTenantAppModule
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantAppModuleRepo) MultiInsert(ctx context.Context, data []*SysTenantAppModule) error {
	if len(data) == 0 {
		return nil
	}
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantAppModule{}).Create(data).Error
	return stores.ErrFmt(err)
}

//func (p TenantAppModuleRepo) MultiUpdate(ctx context.Context, appCode string, moduleCodes []string) error {
//	var datas []*SysTenantAppModule
//	for _, v := range moduleCodes {
//		datas = append(datas, &SysTenantAppModule{
//			AppCodes:    appCode,
//			ModuleCode: v,
//		})
//	}
//	err := p.db.Transaction(func(tx *gorm.DB) error {
//		rm := NewTenantAppModuleRepo(tx)
//		err := rm.DeleteByFilter(ctx, TenantAppModuleFilter{AppCodes: []string{appCode}})
//		if err != nil {
//			return err
//		}
//		if len(datas) != 0 {
//			err = rm.MultiInsert(ctx, datas)
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//	return stores.ErrFmt(err)
//}
