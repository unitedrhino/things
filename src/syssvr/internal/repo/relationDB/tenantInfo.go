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
1. 将TenantInfo全局替换为模型的表名
2. 完善todo
*/

type TenantInfoRepo struct {
	db *gorm.DB
}

func NewTenantInfoRepo(in any) *TenantInfoRepo {
	return &TenantInfoRepo{db: stores.GetCommonConn(in)}
}

type TenantInfoFilter struct {
	ID    int64
	Codes []string
	Code  string
	Name  string
}

func (p TenantInfoRepo) fmtFilter(ctx context.Context, f TenantInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if len(f.Codes) > 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		db = db.Where("code like ?", "%"+f.Code+"%")
	}
	return db
}

func (p TenantInfoRepo) Insert(ctx context.Context, data *SysTenantInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantInfoRepo) FindOneByFilter(ctx context.Context, f TenantInfoFilter) (*SysTenantInfo, error) {
	var result SysTenantInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantInfoRepo) FindByFilter(ctx context.Context, f TenantInfoFilter, page *def.PageInfo) ([]*SysTenantInfo, error) {
	var results []*SysTenantInfo
	db := p.fmtFilter(ctx, f).Model(&SysTenantInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantInfoRepo) CountByFilter(ctx context.Context, f TenantInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantInfoRepo) Update(ctx context.Context, data *SysTenantInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantInfoRepo) DeleteByFilter(ctx context.Context, f TenantInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantInfo{}).Error
	return stores.ErrFmt(err)
}

func (p TenantInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantInfo{}).Error
	return stores.ErrFmt(err)
}
func (p TenantInfoRepo) FindOne(ctx context.Context, id int64) (*SysTenantInfo, error) {
	var result SysTenantInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantInfoRepo) MultiInsert(ctx context.Context, data []*SysTenantInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
