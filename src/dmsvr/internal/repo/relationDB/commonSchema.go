package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type CommonSchemaRepo struct {
	db *gorm.DB
}

type (
	CommonSchemaFilter struct {
		ID          int64
		ProductID   string   //产品id  必填
		Type        int64    //物模型类型 1:property属性 2:event事件 3:action行为
		Tag         int64    //过滤条件: 物模型标签 1:自定义 2:可选 3:必选
		Identifiers []string //过滤标识符列表
	}
)

func NewCommonSchemaRepo(in any) *CommonSchemaRepo {
	return &CommonSchemaRepo{db: stores.GetCommonConn(in)}
}

func (p CommonSchemaRepo) fmtFilter(ctx context.Context, f CommonSchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.Type != 0 {
		db = db.Where("type=?", f.Type)
	}
	if f.Tag != 0 {
		db = db.Where("tag=?", f.Tag)
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

func (p CommonSchemaRepo) FindByFilter(ctx context.Context, f CommonSchemaFilter, page *def.PageInfo) ([]*DmCommonSchema, error) {
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
