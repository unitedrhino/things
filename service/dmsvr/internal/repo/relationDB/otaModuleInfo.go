package relationDB

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将OtaModuleInfo全局替换为模型的表名
2. 完善todo
*/

type OtaModuleInfoRepo struct {
	db *gorm.DB
}

func NewOtaModuleInfoRepo(in any) *OtaModuleInfoRepo {
	return &OtaModuleInfoRepo{db: stores.GetCommonConn(in)}
}

type OtaModuleInfoFilter struct {
	ID        int64
	Code      string
	Name      string
	ProductID string
}

func (p OtaModuleInfoRepo) fmtFilter(ctx context.Context, f OtaModuleInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.Code != "" {
		db = db.Where("code=?", f.Code)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name, "%")
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	return db
}

func (p OtaModuleInfoRepo) Insert(ctx context.Context, data *DmOtaModuleInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OtaModuleInfoRepo) FindOneByFilter(ctx context.Context, f OtaModuleInfoFilter) (*DmOtaModuleInfo, error) {
	var result DmOtaModuleInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaModuleInfoRepo) FindByFilter(ctx context.Context, f OtaModuleInfoFilter, page *stores.PageInfo) ([]*DmOtaModuleInfo, error) {
	var results []*DmOtaModuleInfo
	db := p.fmtFilter(ctx, f).Model(&DmOtaModuleInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaModuleInfoRepo) CountByFilter(ctx context.Context, f OtaModuleInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaModuleInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OtaModuleInfoRepo) Update(ctx context.Context, data *DmOtaModuleInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OtaModuleInfoRepo) DeleteByFilter(ctx context.Context, f OtaModuleInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaModuleInfo{}).Error
	return stores.ErrFmt(err)
}

func (p OtaModuleInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmOtaModuleInfo{}).Error
	return stores.ErrFmt(err)
}
func (p OtaModuleInfoRepo) FindOne(ctx context.Context, id int64) (*DmOtaModuleInfo, error) {
	var result DmOtaModuleInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OtaModuleInfoRepo) MultiInsert(ctx context.Context, data []*DmOtaModuleInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaModuleInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
