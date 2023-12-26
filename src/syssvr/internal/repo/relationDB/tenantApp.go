package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将TenantApp全局替换为模型的表名
2. 完善todo
*/

type TenantAppRepo struct {
	db *gorm.DB
}

func NewTenantAppRepo(in any) *TenantAppRepo {
	return &TenantAppRepo{db: stores.GetCommonConn(in)}
}

type TenantAppFilter struct {
	TenantCode string
	IDs        []int64
	Codes      []string
	//todo 添加过滤字段
}

func (p TenantAppRepo) fmtFilter(ctx context.Context, f TenantAppFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if len(f.IDs) > 0 {
		db = db.Where("id in ?", f.IDs)
	}
	if len(f.Codes) > 0 {
		db = db.Where("code in ?", f.Codes)
	}
	return db
}

func (p TenantAppRepo) Insert(ctx context.Context, data *SysTenantApp) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantAppRepo) FindOneByFilter(ctx context.Context, f TenantAppFilter) (*SysTenantApp, error) {
	var result SysTenantApp
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantAppRepo) FindByFilter(ctx context.Context, f TenantAppFilter, page *def.PageInfo) ([]*SysTenantApp, error) {
	var results []*SysTenantApp
	db := p.fmtFilter(ctx, f).Model(&SysTenantApp{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantAppRepo) CountByFilter(ctx context.Context, f TenantAppFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantApp{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantAppRepo) Update(ctx context.Context, data *SysTenantApp) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantAppRepo) DeleteByFilter(ctx context.Context, f TenantAppFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantApp{}).Error
	return stores.ErrFmt(err)
}

func (p TenantAppRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantApp{}).Error
	return stores.ErrFmt(err)
}
func (p TenantAppRepo) FindOne(ctx context.Context, id int64) (*SysTenantApp, error) {
	var result SysTenantApp
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantAppRepo) MultiInsert(ctx context.Context, data []*SysTenantApp) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantApp{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p TenantAppRepo) MultiUpdate(ctx context.Context, tenantCode string, appCodes []string) error {
	var datas []*SysTenantApp
	for _, v := range appCodes {
		datas = append(datas, &SysTenantApp{
			AppCode:    v,
			TenantCode: stores.TenantCode(tenantCode),
		})
	}
	ctxs.GetUserCtx(ctx).AllTenant = true //只有这样才能改为其他租户
	defer func() {
		ctxs.GetUserCtx(ctx).AllTenant = false
	}()
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewTenantAppRepo(tx)
		err := rm.DeleteByFilter(ctx, TenantAppFilter{TenantCode: tenantCode})
		if err != nil {
			return err
		}
		if len(datas) != 0 {
			err = rm.MultiInsert(ctx, datas)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}
