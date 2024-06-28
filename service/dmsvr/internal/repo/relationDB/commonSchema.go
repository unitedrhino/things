package relationDB

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
)

type CommonSchemaRepo struct {
	db *gorm.DB
}

type (
	CommonSchemaFilter struct {
		ID                int64
		Type              int64 //物模型类型 1:property属性 2:event事件 3:action行为
		Types             []int64
		Identifiers       []string //过滤标识符列表
		Name              string
		IsCanSceneLinkage int64
		FuncGroup         int64
		UserPerm          int64
		PropertyMode      string
	}
)

func NewCommonSchemaRepo(in any) *CommonSchemaRepo {
	return &CommonSchemaRepo{db: stores.GetCommonConn(in)}
}

func (p CommonSchemaRepo) fmtFilter(ctx context.Context, f CommonSchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.IsCanSceneLinkage != 0 {
		db = db.Where("is_can_scene_linkage = ?", f.IsCanSceneLinkage)
	}
	if f.FuncGroup != 0 {
		db = db.Where("func_group = ?", f.FuncGroup)
	}
	if f.UserPerm != 0 {
		db = db.Where("user_auth = ?", f.UserPerm)
	}
	if f.PropertyMode != "" {
		db = db.Where("JSON_CONTAINS(affordance, JSON_OBJECT('mode',?))",
			f.PropertyMode)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.Type != 0 {
		db = db.Where("type=?", f.Type)
	}
	if len(f.Types) != 0 {
		db = db.Where("type in ?", f.Types)
	}
	if len(f.Identifiers) != 0 {
		db = db.Where("identifier in ?", f.Identifiers)
	}
	return db
}
func (p CommonSchemaRepo) Insert(ctx context.Context, data *DmCommonSchema) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p CommonSchemaRepo) FindOneByFilter(ctx context.Context, f CommonSchemaFilter) (*DmCommonSchema, error) {
	var result DmCommonSchema
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p CommonSchemaRepo) Update(ctx context.Context, data *DmCommonSchema) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p CommonSchemaRepo) DeleteByFilter(ctx context.Context, f CommonSchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmCommonSchema{}).Error
	return stores.ErrFmt(err)
}

func (p CommonSchemaRepo) FindByFilter(ctx context.Context, f CommonSchemaFilter, page *stores.PageInfo) ([]*DmCommonSchema, error) {
	var results []*DmCommonSchema
	db := p.fmtFilter(ctx, f).Model(&DmCommonSchema{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p CommonSchemaRepo) CountByFilter(ctx context.Context, f CommonSchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmCommonSchema{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
