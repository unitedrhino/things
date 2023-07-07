package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"gorm.io/gorm"
)

type MenuInfoRepo struct {
	db *gorm.DB
}

func NewMenuInfoRepo(in any) *MenuInfoRepo {
	return &MenuInfoRepo{db: store.GetCommonConn(in)}
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

func (g MenuInfoRepo) Insert(ctx context.Context, data *SysMenuInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return store.ErrFmt(result.Error)
}

func (g MenuInfoRepo) FindOneByFilter(ctx context.Context, f MenuInfoFilter) (*SysMenuInfo, error) {
	var result SysMenuInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}
func (p MenuInfoRepo) FindByFilter(ctx context.Context, f MenuInfoFilter, page *def.PageInfo) ([]*SysMenuInfo, error) {
	var results []*SysMenuInfo
	db := p.fmtFilter(ctx, f).Model(&SysMenuInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return results, nil
}

func (p MenuInfoRepo) CountByFilter(ctx context.Context, f MenuInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysMenuInfo{})
	err = db.Count(&size).Error
	return size, store.ErrFmt(err)
}

func (g MenuInfoRepo) Update(ctx context.Context, data *SysMenuInfo) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return store.ErrFmt(err)
}

func (g MenuInfoRepo) DeleteByFilter(ctx context.Context, f MenuInfoFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysMenuInfo{}).Error
	return store.ErrFmt(err)
}

func (g MenuInfoRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&SysMenuInfo{}).Error
	return store.ErrFmt(err)
}
func (g MenuInfoRepo) FindOne(ctx context.Context, id int64) (*SysMenuInfo, error) {
	var result SysMenuInfo
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}
