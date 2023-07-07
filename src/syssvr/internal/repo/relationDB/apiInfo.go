package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"gorm.io/gorm"
)

type ApiInfoRepo struct {
	db *gorm.DB
}

func NewApiInfoRepo(in any) *ApiInfoRepo {
	return &ApiInfoRepo{db: store.GetCommonConn(in)}
}

type ApiInfoFilter struct {
	Route  string
	Method int64
	Group  string
	Name   string
}

func (p ApiInfoRepo) fmtFilter(ctx context.Context, f ApiInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Route != "" {
		db = db.Where("`route` like ?", "%"+f.Route+"%")
	}
	if f.Method != 0 {
		db = db.Where("`method` = ?", f.Method)
	}
	if f.Group != "" {
		db = db.Where("`group` like ?", "%"+f.Group+"%")
	}
	if f.Name != "" {
		db = db.Where("`name` like ?", "%"+f.Name+"%")
	}
	return db
}

func (g ApiInfoRepo) Insert(ctx context.Context, data *SysApiInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return store.ErrFmt(result.Error)
}

func (g ApiInfoRepo) FindOneByFilter(ctx context.Context, f ApiInfoFilter) (*SysApiInfo, error) {
	var result SysApiInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}

func (p ApiInfoRepo) FindByFilter(ctx context.Context, f ApiInfoFilter, page *def.PageInfo) ([]*SysApiInfo, error) {
	var results []*SysApiInfo
	db := p.fmtFilter(ctx, f).Model(&SysApiInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return results, nil
}

func (p ApiInfoRepo) CountByFilter(ctx context.Context, f ApiInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysApiInfo{})
	err = db.Count(&size).Error
	return size, store.ErrFmt(err)
}

func (g ApiInfoRepo) Update(ctx context.Context, data *SysApiInfo) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return store.ErrFmt(err)
}

func (g ApiInfoRepo) DeleteByFilter(ctx context.Context, f ApiInfoFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysApiInfo{}).Error
	return store.ErrFmt(err)
}
func (g ApiInfoRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&SysApiInfo{}).Error
	return store.ErrFmt(err)
}

func (g ApiInfoRepo) FindOne(ctx context.Context, id int64) (*SysApiInfo, error) {
	var result SysApiInfo
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}
