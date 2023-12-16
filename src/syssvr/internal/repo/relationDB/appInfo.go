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
1. 将example全局替换为模型的表名
2. 完善todo
*/

type AppInfoRepo struct {
	db *gorm.DB
}

func NewAppInfoRepo(in any) *AppInfoRepo {
	return &AppInfoRepo{db: stores.GetCommonConn(in)}
}

type AppInfoFilter struct {
	ID    int64
	Codes []string
	Code  string
	Name  string
}

func (p AppInfoRepo) fmtFilter(ctx context.Context, f AppInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
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

func (p AppInfoRepo) Insert(ctx context.Context, data *SysAppInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AppInfoRepo) FindOneByFilter(ctx context.Context, f AppInfoFilter) (*SysAppInfo, error) {
	var result SysAppInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AppInfoRepo) FindByFilter(ctx context.Context, f AppInfoFilter, page *def.PageInfo) ([]*SysAppInfo, error) {
	var results []*SysAppInfo
	db := p.fmtFilter(ctx, f).Model(&SysAppInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AppInfoRepo) CountByFilter(ctx context.Context, f AppInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysAppInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AppInfoRepo) Update(ctx context.Context, data *SysAppInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AppInfoRepo) DeleteByFilter(ctx context.Context, f AppInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysAppInfo{}).Error
	return stores.ErrFmt(err)
}

func (p AppInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysAppInfo{}).Error
	return stores.ErrFmt(err)
}
func (p AppInfoRepo) FindOne(ctx context.Context, id int64) (*SysAppInfo, error) {
	var result SysAppInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AppInfoRepo) MultiInsert(ctx context.Context, data []*SysAppInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysAppInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
