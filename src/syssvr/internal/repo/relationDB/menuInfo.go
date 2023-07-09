package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type MenuInfoRepo struct {
	db *gorm.DB
}

func NewMenuInfoRepo(in any) *MenuInfoRepo {
	return &MenuInfoRepo{db: stores.GetCommonConn(in)}
}

type MenuInfoFilter struct {
	Role    int64
	Name    string
	Path    string
	MenuIds []int64
}

func (p MenuInfoRepo) fmtFilter(ctx context.Context, f MenuInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("`name` like ?", "%"+f.Name+"%")
	}
	if f.Path != "" {
		db = db.Where("`path` like ?", "%"+f.Path+"%")
	}
	if len(f.MenuIds) != 0 {
		db = db.Where("`id` in ?", f.MenuIds)
	}
	return db
}

func (p MenuInfoRepo) Insert(ctx context.Context, data *SysMenuInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p MenuInfoRepo) FindOneByFilter(ctx context.Context, f MenuInfoFilter) (*SysMenuInfo, error) {
	var result SysMenuInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p MenuInfoRepo) FindByFilter(ctx context.Context, f MenuInfoFilter, page *def.PageInfo) ([]*SysMenuInfo, error) {
	var results []*SysMenuInfo
	db := p.fmtFilter(ctx, f).Model(&SysMenuInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p MenuInfoRepo) CountByFilter(ctx context.Context, f MenuInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysMenuInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p MenuInfoRepo) Update(ctx context.Context, data *SysMenuInfo) error {
	err := p.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p MenuInfoRepo) DeleteByFilter(ctx context.Context, f MenuInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysMenuInfo{}).Error
	return stores.ErrFmt(err)
}

func (p MenuInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("`id` = ?", id).Delete(&SysMenuInfo{}).Error
	return stores.ErrFmt(err)
}
func (p MenuInfoRepo) FindOne(ctx context.Context, id int64) (*SysMenuInfo, error) {
	var result SysMenuInfo
	err := p.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
