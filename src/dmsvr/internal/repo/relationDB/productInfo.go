package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type ProductInfoRepo struct {
	db *gorm.DB
}

type ProductFilter struct {
	DeviceType   int64
	ProductName  string
	ProductIDs   []string
	ProductNames []string
	Tags         map[string]string
}

func NewProductInfoRepo(in any) *ProductInfoRepo {
	return &ProductInfoRepo{db: stores.GetCommonConn(in)}
}

func (p ProductInfoRepo) fmtFilter(ctx context.Context, f ProductFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.DeviceType != 0 {
		db = db.Where("device_type=?", f.DeviceType)
	}
	if f.ProductName != "" {
		db = db.Where("product_name like ?", "%"+f.ProductName+"%")
	}
	if len(f.ProductIDs) != 0 {
		db = db.Where("product_id in ?", f.ProductIDs)
	}
	if len(f.ProductNames) != 0 {
		db = db.Where("product_name in ?", f.ProductNames)
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = db.Where("JSON_CONTAINS(tags, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	return db
}

func (p ProductInfoRepo) Insert(ctx context.Context, data *DmProductInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductInfoRepo) FindOneByFilter(ctx context.Context, f ProductFilter) (*DmProductInfo, error) {
	var result DmProductInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductInfoRepo) Update(ctx context.Context, data *DmProductInfo) error {
	err := p.db.WithContext(ctx).Where("product_id = ?", data.ProductID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) DeleteByFilter(ctx context.Context, f ProductFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProductInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ProductInfoRepo) FindByFilter(ctx context.Context, f ProductFilter, page *def.PageInfo) ([]*DmProductInfo, error) {
	var results []*DmProductInfo
	db := p.fmtFilter(ctx, f).Model(&DmProductInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProductInfoRepo) CountByFilter(ctx context.Context, f ProductFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProductInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
