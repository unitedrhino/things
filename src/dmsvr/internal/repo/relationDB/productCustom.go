package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/store"
	"gorm.io/gorm"
)

type ProductCustomRepo struct {
	db *gorm.DB
}

func NewProductCustomRepo(in any) *ProductCustomRepo {
	return &ProductCustomRepo{db: store.GetCommonConn(in)}
}

func (p ProductCustomRepo) Insert(ctx context.Context, data *DmProductCustom) error {
	result := p.db.WithContext(ctx).Create(data)
	return store.ErrFmt(result.Error)
}

func (p ProductCustomRepo) FindOneByProductID(ctx context.Context, productID string) (*DmProductCustom, error) {
	var result DmProductCustom
	err := p.db.WithContext(ctx).Where("productID = ?", productID).First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}

func (p ProductCustomRepo) Update(ctx context.Context, data *DmProductCustom) error {
	err := p.db.WithContext(ctx).Save(data).Error
	return store.ErrFmt(err)
}
