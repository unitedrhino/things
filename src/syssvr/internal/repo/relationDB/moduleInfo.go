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
1. 将ModuleInfo全局替换为模型的表名
2. 完善todo
*/

type ModuleInfoRepo struct {
	db *gorm.DB
}

func NewModuleInfoRepo(in any) *ModuleInfoRepo {
	return &ModuleInfoRepo{db: stores.GetCommonConn(in)}
}

type ModuleInfoFilter struct {
	ID        int64
	Codes     []string
	Code      string
	Name      string
	WithApis  bool
	WithMenus bool
}

func (p ModuleInfoRepo) fmtFilter(ctx context.Context, f ModuleInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithApis {
		db = db.Preload("Apis")
	}
	if f.WithMenus {
		db = db.Preload("Menus")
	}
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if len(f.Codes) > 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		db = db.Where("code like ?", "%"+f.Code+"%")
	}
	return db
}

func (p ModuleInfoRepo) Insert(ctx context.Context, data *SysModuleInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ModuleInfoRepo) FindOneByFilter(ctx context.Context, f ModuleInfoFilter) (*SysModuleInfo, error) {
	var result SysModuleInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ModuleInfoRepo) FindByFilter(ctx context.Context, f ModuleInfoFilter, page *def.PageInfo) ([]*SysModuleInfo, error) {
	var results []*SysModuleInfo
	db := p.fmtFilter(ctx, f).Model(&SysModuleInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ModuleInfoRepo) CountByFilter(ctx context.Context, f ModuleInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysModuleInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ModuleInfoRepo) Update(ctx context.Context, data *SysModuleInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ModuleInfoRepo) DeleteByFilter(ctx context.Context, f ModuleInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysModuleInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ModuleInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysModuleInfo{}).Error
	return stores.ErrFmt(err)
}
func (p ModuleInfoRepo) FindOne(ctx context.Context, id int64) (*SysModuleInfo, error) {
	var result SysModuleInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ModuleInfoRepo) MultiInsert(ctx context.Context, data []*SysModuleInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysModuleInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
