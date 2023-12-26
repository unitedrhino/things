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
1. 将AppModule全局替换为模型的表名
2. 完善todo
*/

type AppModuleRepo struct {
	db *gorm.DB
}

func NewAppModuleRepo(in any) *AppModuleRepo {
	return &AppModuleRepo{db: stores.GetCommonConn(in)}
}

type AppModuleFilter struct {
	AppCode []string
}

func (p AppModuleRepo) fmtFilter(ctx context.Context, f AppModuleFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.AppCode) != 0 {
		db = db.Where("app_code in ?", f.AppCode)
	}
	return db
}

func (p AppModuleRepo) Insert(ctx context.Context, data *SysAppModule) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AppModuleRepo) FindOneByFilter(ctx context.Context, f AppModuleFilter) (*SysAppModule, error) {
	var result SysAppModule
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AppModuleRepo) FindByFilter(ctx context.Context, f AppModuleFilter, page *def.PageInfo) ([]*SysAppModule, error) {
	var results []*SysAppModule
	db := p.fmtFilter(ctx, f).Model(&SysAppModule{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AppModuleRepo) CountByFilter(ctx context.Context, f AppModuleFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysAppModule{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AppModuleRepo) Update(ctx context.Context, data *SysAppModule) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AppModuleRepo) DeleteByFilter(ctx context.Context, f AppModuleFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysAppModule{}).Error
	return stores.ErrFmt(err)
}

func (p AppModuleRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysAppModule{}).Error
	return stores.ErrFmt(err)
}
func (p AppModuleRepo) FindOne(ctx context.Context, id int64) (*SysAppModule, error) {
	var result SysAppModule
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AppModuleRepo) MultiInsert(ctx context.Context, data []*SysAppModule) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysAppModule{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p AppModuleRepo) MultiUpdate(ctx context.Context, appCode string, moduleCodes []string) error {
	var datas []*SysAppModule
	for _, v := range moduleCodes {
		datas = append(datas, &SysAppModule{
			AppCode:    appCode,
			ModuleCode: v,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewAppModuleRepo(tx)
		err := rm.DeleteByFilter(ctx, AppModuleFilter{AppCode: []string{appCode}})
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
