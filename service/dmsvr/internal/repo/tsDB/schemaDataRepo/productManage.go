package schemaDataRepo

import (
	"context"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/stores"
)

func (d *DeviceDataRepo) DeleteProduct(ctx context.Context, t *schema.Model, productID string) error {
	var typ = map[schema.DataType]struct{}{}
	for _, p := range t.Property {
		if _, ok := typ[p.Define.Type]; ok {
			continue
		}
		typ[p.Define.Type] = struct{}{}
		err := d.db.WithContext(ctx).Table(getTableName(p.Define)).
			Where("product_id = ?", productID).Delete(&Property{}).Error
		if err != nil {
			return stores.ErrFmt(err)
		}
	}
	return nil
}

func (d *DeviceDataRepo) InitProduct(ctx context.Context, t *schema.Model, productID string) error {
	return nil
}

func (d *DeviceDataRepo) CreateProperty(ctx context.Context, p *schema.Property, productID string) error {
	return nil
}
func (d *DeviceDataRepo) DeleteProperty(ctx context.Context, p *schema.Property, productID string, identifier string) error {
	err := d.db.WithContext(ctx).Table(getTableName(p.Define)).Where("product_id = ? AND identifier = ?", productID, identifier).Delete(&Property{}).Error
	return stores.ErrFmt(err)
}

func (d *DeviceDataRepo) UpdateProperty(
	ctx context.Context,
	oldP *schema.Property,
	newP *schema.Property,
	productID string) error {
	return nil
}
