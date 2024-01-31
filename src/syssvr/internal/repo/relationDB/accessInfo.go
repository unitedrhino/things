package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type AccessRepo struct {
	db *gorm.DB
}

func NewAccessRepo(in any) *AccessRepo {
	return &AccessRepo{db: stores.GetCommonConn(in)}
}

type AccessFilter struct {
	Name       string
	Code       string
	Codes      []string
	IsNeedAuth int64
	Group      string
	WithApis   bool
}

func (p AccessRepo) fmtFilter(ctx context.Context, f AccessFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithApis {
		db = db.Preload("Apis")
	}
	if f.IsNeedAuth != 0 {
		db = db.Where("is_need_auth =?", f.IsNeedAuth)
	}
	if len(f.Codes) != 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Code != "" {
		db = db.Where("code like ? ", "%"+f.Code+"%")
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Group != "" {
		db = db.Where("group like ?", "%"+f.Group+"%")
	}
	return db
}

func (p AccessRepo) Insert(ctx context.Context, data *SysAccessInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AccessRepo) FindOneByFilter(ctx context.Context, f AccessFilter) (*SysAccessInfo, error) {
	var result SysAccessInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p AccessRepo) FindByFilter(ctx context.Context, f AccessFilter, page *def.PageInfo) ([]*SysAccessInfo, error) {
	var results []*SysAccessInfo
	db := p.fmtFilter(ctx, f).Model(&SysAccessInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AccessRepo) CountByFilter(ctx context.Context, f AccessFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysAccessInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AccessRepo) Update(ctx context.Context, data *SysAccessInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AccessRepo) DeleteByFilter(ctx context.Context, f AccessFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysAccessInfo{}).Error
	return stores.ErrFmt(err)
}
func (p AccessRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysAccessInfo{}).Error
	return stores.ErrFmt(err)
}

func (p AccessRepo) FindOne(ctx context.Context, id int64) (*SysAccessInfo, error) {
	var result SysAccessInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
