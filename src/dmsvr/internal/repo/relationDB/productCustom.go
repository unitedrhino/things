package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type ProductCustomRepo struct {
	db *gorm.DB
}

func NewProductCustomRepo(in any) *ProductCustomRepo {
	return &ProductCustomRepo{db: stores.GetCommonConn(in)}
}

func (p ProductCustomRepo) Insert(ctx context.Context, data *DmProductCustom) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProductCustomRepo) FindOneByProductID(ctx context.Context, productID string) (*DmProductCustom, error) {
	var result DmProductCustom
	err := p.db.WithContext(ctx).Where("product_id = ?", productID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductCustomRepo) Update(ctx context.Context, data *DmProductCustom) error {
	err := p.db.WithContext(ctx).Save(data).Error
	return stores.ErrFmt(err)
}
