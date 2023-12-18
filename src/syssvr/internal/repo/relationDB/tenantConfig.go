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
1. 将TenantConfig全局替换为模型的表名
2. 完善todo
*/

type TenantConfigRepo struct {
	db *gorm.DB
}

func NewTenantConfigRepo(in any) *TenantConfigRepo {
	return &TenantConfigRepo{db: stores.GetCommonConn(in)}
}

type TenantConfigFilter struct {
	//todo 添加过滤字段
}

func (p TenantConfigRepo) fmtFilter(ctx context.Context, f TenantConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p TenantConfigRepo) Insert(ctx context.Context, data *SysTenantConfig) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantConfigRepo) FindOneByFilter(ctx context.Context, f TenantConfigFilter) (*SysTenantConfig, error) {
	var result SysTenantConfig
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantConfigRepo) FindByFilter(ctx context.Context, f TenantConfigFilter, page *def.PageInfo) ([]*SysTenantConfig, error) {
	var results []*SysTenantConfig
	db := p.fmtFilter(ctx, f).Model(&SysTenantConfig{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantConfigRepo) CountByFilter(ctx context.Context, f TenantConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantConfig{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantConfigRepo) Update(ctx context.Context, data *SysTenantConfig) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantConfigRepo) DeleteByFilter(ctx context.Context, f TenantConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantConfig{}).Error
	return stores.ErrFmt(err)
}

func (p TenantConfigRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantConfig{}).Error
	return stores.ErrFmt(err)
}
func (p TenantConfigRepo) FindOne(ctx context.Context) (*SysTenantConfig, error) {
	var result SysTenantConfig
	err := p.db.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantConfigRepo) MultiInsert(ctx context.Context, data []*SysTenantConfig) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantConfig{}).Create(data).Error
	return stores.ErrFmt(err)
}
