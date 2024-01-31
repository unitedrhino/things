package relationDB

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/stores"
	"gorm.io/gorm"
)

type ProductRemoteConfigRepo struct {
	db *gorm.DB
}

type (
	RemoteConfigFilter struct {
		ProductID string
	}
)

func NewProductRemoteConfigRepo(in any) *ProductRemoteConfigRepo {
	return &ProductRemoteConfigRepo{db: stores.GetCommonConn(in)}
}

func (p ProductRemoteConfigRepo) fmtFilter(ctx context.Context, f RemoteConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = db.Order("created_time desc")
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	return db
}

func (p ProductRemoteConfigRepo) Insert(ctx context.Context, data *DmProductRemoteConfig) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductRemoteConfigRepo) FindByFilter(ctx context.Context, f RemoteConfigFilter, page *def.PageInfo) ([]*DmProductRemoteConfig, error) {
	var results []*DmProductRemoteConfig
	db := p.fmtFilter(ctx, f).Model(&DmProductRemoteConfig{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductRemoteConfigRepo) FindOneByFilter(ctx context.Context, f RemoteConfigFilter) (*DmProductRemoteConfig, error) {
	var result DmProductRemoteConfig
	db := p.fmtFilter(ctx, f)
	err := db.Last(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductRemoteConfigRepo) Update(ctx context.Context, data *DmProductRemoteConfig) error {
	err := p.db.WithContext(ctx).Where("product_id = ?", data.ProductID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductRemoteConfigRepo) DeleteByFilter(ctx context.Context, f RemoteConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductRemoteConfig{}).Error
	return stores.ErrFmt(err)
}
func (p ProductRemoteConfigRepo) CountByFilter(ctx context.Context, f RemoteConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductRemoteConfig{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
