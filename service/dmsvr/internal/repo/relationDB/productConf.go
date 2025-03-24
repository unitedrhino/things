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
1. 将ProductConfig全局替换为模型的表名
2. 完善todo
*/

type ProductConfigRepo struct {
	db *gorm.DB
}

func NewProductConfigRepo(in any) *ProductConfigRepo {
	return &ProductConfigRepo{db: stores.GetCommonConn(in)}
}

type ProductConfigFilter struct {
	//todo 添加过滤字段
}

func (p ProductConfigRepo) fmtFilter(ctx context.Context, f ProductConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p ProductConfigRepo) Insert(ctx context.Context, data *DmProductConfig) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductConfigRepo) FindOneByFilter(ctx context.Context, f ProductConfigFilter) (*DmProductConfig, error) {
	var result DmProductConfig
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProductConfigRepo) FindByFilter(ctx context.Context, f ProductConfigFilter, page *stores.PageInfo) ([]*DmProductConfig, error) {
	var results []*DmProductConfig
	db := p.fmtFilter(ctx, f).Model(&DmProductConfig{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductConfigRepo) CountByFilter(ctx context.Context, f ProductConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductConfig{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProductConfigRepo) Update(ctx context.Context, data *DmProductConfig) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductConfigRepo) DeleteByFilter(ctx context.Context, f ProductConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductConfig{}).Error
	return stores.ErrFmt(err)
}

func (p ProductConfigRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProductConfig{}).Error
	return stores.ErrFmt(err)
}
func (p ProductConfigRepo) FindOne(ctx context.Context, productID string) (*DmProductConfig, error) {
	var result DmProductConfig
	err := p.db.WithContext(ctx).Where("product_id = ?", productID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProductConfigRepo) MultiInsert(ctx context.Context, data []*DmProductConfig) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProductConfig{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d ProductConfigRepo) UpdateWithField(ctx context.Context, f ProductConfigFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProductConfig{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
