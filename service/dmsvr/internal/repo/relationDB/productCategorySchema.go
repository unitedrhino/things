package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将ProductCategorySchema全局替换为模型的表名
2. 完善todo
*/

type ProductCategorySchemaRepo struct {
	db *gorm.DB
}

func NewProductCategorySchemaRepo(in any) *ProductCategorySchemaRepo {
	return &ProductCategorySchemaRepo{db: stores.GetCommonConn(in)}
}

type ProductCategorySchemaFilter struct {
	ProductCategoryID  int64
	ProductCategoryIDs []int64
	Identifiers        []string
}

func (p ProductCategorySchemaRepo) fmtFilter(ctx context.Context, f ProductCategorySchemaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ProductCategoryID != 0 {
		db = db.Where("product_category_id=?", f.ProductCategoryID)
	}
	if len(f.ProductCategoryIDs) != 0 {
		db = db.Where("product_category_id in ?", f.ProductCategoryIDs)
	}
	if len(f.Identifiers) != 0 {
		db = db.Where("identifier in ?", f.Identifiers)
	}
	return db
}

func (p ProductCategorySchemaRepo) Insert(ctx context.Context, data *DmProductCategorySchema) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductCategorySchemaRepo) FindOneByFilter(ctx context.Context, f ProductCategorySchemaFilter) (*DmProductCategorySchema, error) {
	var result DmProductCategorySchema
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProductCategorySchemaRepo) FindByFilter(ctx context.Context, f ProductCategorySchemaFilter, page *def.PageInfo) ([]*DmProductCategorySchema, error) {
	var results []*DmProductCategorySchema
	db := p.fmtFilter(ctx, f).Model(&DmProductCategorySchema{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductCategorySchemaRepo) CountByFilter(ctx context.Context, f ProductCategorySchemaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductCategorySchema{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProductCategorySchemaRepo) Update(ctx context.Context, data *DmProductCategorySchema) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductCategorySchemaRepo) DeleteByFilter(ctx context.Context, f ProductCategorySchemaFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductCategorySchema{}).Error
	return stores.ErrFmt(err)
}

func (p ProductCategorySchemaRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProductCategorySchema{}).Error
	return stores.ErrFmt(err)
}
func (p ProductCategorySchemaRepo) FindOne(ctx context.Context, id int64) (*DmProductCategorySchema, error) {
	var result DmProductCategorySchema
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProductCategorySchemaRepo) MultiInsert(ctx context.Context, data []*DmProductCategorySchema) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProductCategorySchema{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p ProductCategorySchemaRepo) MultiUpdate(ctx context.Context, productCategoryID int64, identifiers []string) error {
	var insertDatas []*DmProductCategorySchema
	for _, v := range identifiers {
		insertDatas = append(insertDatas, &DmProductCategorySchema{
			ProductCategoryID: productCategoryID,
			Identifier:        v,
		})
	}

	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewProductCategorySchemaRepo(tx)
		err := rm.DeleteByFilter(ctx, ProductCategorySchemaFilter{ProductCategoryID: productCategoryID})
		if err != nil {
			return err
		}
		if len(identifiers) != 0 {
			err = rm.MultiInsert(ctx, insertDatas)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}
