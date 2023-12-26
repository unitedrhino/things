package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type ApiInfoRepo struct {
	db *gorm.DB
}

func NewApiInfoRepo(in any) *ApiInfoRepo {
	return &ApiInfoRepo{db: stores.GetCommonConn(in)}
}

type ApiInfoFilter struct {
	ApiIDs     []int64
	Route      string
	Method     string
	Group      string
	Name       string
	ModuleCode string
	IsNeedAuth int64
}

func (p ApiInfoRepo) fmtFilter(ctx context.Context, f ApiInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ApiIDs != nil {
		db = db.Where("id in ?", f.ApiIDs)
	}
	if f.IsNeedAuth != 0 {
		db = db.Where("is_need_auth =?", f.IsNeedAuth)
	}
	if f.Route != "" {
		db = db.Where("route like ?", "%"+f.Route+"%")
	}
	if f.ModuleCode != "" {
		db = db.Where("module_code =?", f.ModuleCode)
	}
	if f.Method != "" {
		db = db.Where("method = ?", f.Method)
	}
	if f.Group != "" {
		db = db.Where("group like ?", "%"+f.Group+"%")
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	return db
}

func (p ApiInfoRepo) Insert(ctx context.Context, data *SysModuleApi) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ApiInfoRepo) FindOneByFilter(ctx context.Context, f ApiInfoFilter) (*SysModuleApi, error) {
	var result SysModuleApi
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ApiInfoRepo) FindByFilter(ctx context.Context, f ApiInfoFilter, page *def.PageInfo) ([]*SysModuleApi, error) {
	var results []*SysModuleApi
	db := p.fmtFilter(ctx, f).Model(&SysModuleApi{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ApiInfoRepo) CountByFilter(ctx context.Context, f ApiInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysModuleApi{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ApiInfoRepo) Update(ctx context.Context, data *SysModuleApi) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ApiInfoRepo) DeleteByFilter(ctx context.Context, f ApiInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysModuleApi{}).Error
	return stores.ErrFmt(err)
}
func (p ApiInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysModuleApi{}).Error
	return stores.ErrFmt(err)
}

func (p ApiInfoRepo) FindOne(ctx context.Context, id int64) (*SysModuleApi, error) {
	var result SysModuleApi
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
