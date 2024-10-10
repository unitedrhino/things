package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type ProductCategoryRepo struct {
	db *gorm.DB
}

func NewProductCategoryRepo(in any) *ProductCategoryRepo {
	return &ProductCategoryRepo{db: stores.GetCommonConn(in)}
}

type ProductCategoryFilter struct {
	Name       string
	IDPath     string
	ParentID   int64
	ID         int64
	IDs        []int64
	ProductIDs []string
}

func (p ProductCategoryRepo) fmtFilter(ctx context.Context, f ProductCategoryFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if len(f.ProductIDs) > 0 {
		subQuery := p.db.Model(&DmProductInfo{}).Select("category_id").Where("product_id in ?", f.ProductIDs)
		db = db.Where("id in (?)",
			subQuery)
	}
	if len(f.IDs) != 0 {
		db = db.Where("id in ?", f.IDs)
	}
	if f.ParentID != 0 {
		db = db.Where("parent_id = ?", f.ParentID)
	}
	if f.IDPath != "" {
		db = db.Where("id_path like ?", f.IDPath+"%")
	}
	return db
}

func (p ProductCategoryRepo) Insert(ctx context.Context, data *DmProductCategory) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductCategoryRepo) FindOneByFilter(ctx context.Context, f ProductCategoryFilter) (*DmProductCategory, error) {
	var result DmProductCategory
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductCategoryRepo) FindByFilter(ctx context.Context, f ProductCategoryFilter, page *stores.PageInfo) ([]*DmProductCategory, error) {
	var results []*DmProductCategory
	db := p.fmtFilter(ctx, f).Model(&DmProductCategory{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductCategoryRepo) CountByFilter(ctx context.Context, f ProductCategoryFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductCategory{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProductCategoryRepo) Update(ctx context.Context, data *DmProductCategory) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d ProductCategoryRepo) UpdateWithField(ctx context.Context, f ProductCategoryFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProductCategory{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p ProductCategoryRepo) DeleteByFilter(ctx context.Context, f ProductCategoryFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductCategory{}).Error
	return stores.ErrFmt(err)
}

func (p ProductCategoryRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProductCategory{}).Error
	return stores.ErrFmt(err)
}
func (p ProductCategoryRepo) FindOne(ctx context.Context, id int64) (*DmProductCategory, error) {
	var result DmProductCategory
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProductCategoryRepo) MultiInsert(ctx context.Context, data []*DmProductCategory) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProductCategory{}).Create(data).Error
	return stores.ErrFmt(err)
}
