package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type OtaModuleInfoRepo struct {
	db *gorm.DB
}

func NewOtaModuleInfoRepo(in any) *OtaModuleInfoRepo {
	return &OtaModuleInfoRepo{db: stores.GetCommonConn(in)}
}

type OtaModuleInfoFilter struct {
	DeviceName string
	ProductId  string
	ModuleName string
}

func (p OtaModuleInfoRepo) fmtFilter(ctx context.Context, f OtaModuleInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ModuleName != "" {
		db = db.Where("module_name = ", f.ModuleName)
	}
	if f.ProductId != "" {
		db = db.Where("product_id=", f.ProductId)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name=", f.DeviceName)
	}
	return db
}

func (p OtaModuleInfoRepo) Insert(ctx context.Context, data *DmOtaModule) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OtaModuleInfoRepo) FindOneByFilter(ctx context.Context, f OtaModuleInfoFilter) (*DmOtaModule, error) {
	var result DmOtaModule
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaModuleInfoRepo) FindByFilter(ctx context.Context, f OtaModuleInfoFilter, page *def.PageInfo) ([]*DmOtaModule, error) {
	var results []*DmOtaModule
	db := p.fmtFilter(ctx, f).Model(&DmOtaModule{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaModuleInfoRepo) CountByFilter(ctx context.Context, f OtaModuleInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaModule{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OtaModuleInfoRepo) Update(ctx context.Context, data *DmOtaModule) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OtaModuleInfoRepo) DeleteByFilter(ctx context.Context, f OtaModuleInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaModule{}).Error
	return stores.ErrFmt(err)
}

func (p OtaModuleInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmOtaModule{}).Error
	return stores.ErrFmt(err)
}
func (p OtaModuleInfoRepo) FindOne(ctx context.Context, id int64) (*DmOtaModule, error) {
	var result DmOtaModule
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OtaModuleInfoRepo) MultiInsert(ctx context.Context, data []*DmOtaModule) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaModule{}).Create(data).Error
	return stores.ErrFmt(err)
}
