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
1. 将example全局替换为模型的表名
2. 完善todo
*/

type ProductIDRepo struct {
	db *gorm.DB
}

func NewProductIDRepo(in any) *ProductIDRepo {
	return &ProductIDRepo{db: stores.GetCommonConn(in)}
}

type ProductIDFilter struct {
	//todo 添加过滤字段
}

func (p ProductIDRepo) fmtFilter(ctx context.Context, f ProductIDFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p ProductIDRepo) Insert(ctx context.Context, data *DmProductID) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}
func (p ProductIDRepo) GenID(ctx context.Context) (int64, error) {
	po := DmProductID{}
	err := p.Insert(ctx, &po)
	if err != nil {
		return 0, err
	}
	return po.ID, err
}

func (p ProductIDRepo) FindOneByFilter(ctx context.Context, f ProductIDFilter) (*DmProductID, error) {
	var result DmProductID
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProductIDRepo) FindByFilter(ctx context.Context, f ProductIDFilter, page *def.PageInfo) ([]*DmProductID, error) {
	var results []*DmProductID
	db := p.fmtFilter(ctx, f).Model(&DmProductID{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductIDRepo) CountByFilter(ctx context.Context, f ProductIDFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductID{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProductIDRepo) Update(ctx context.Context, data *DmProductID) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductIDRepo) DeleteByFilter(ctx context.Context, f ProductIDFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductID{}).Error
	return stores.ErrFmt(err)
}

func (p ProductIDRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProductID{}).Error
	return stores.ErrFmt(err)
}
func (p ProductIDRepo) FindOne(ctx context.Context, id int64) (*DmProductID, error) {
	var result DmProductID
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProductIDRepo) MultiInsert(ctx context.Context, data []*DmProductID) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProductID{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d ProductIDRepo) UpdateWithField(ctx context.Context, f ProductIDFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProductID{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
