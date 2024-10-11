package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/shadow"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShadowRepo struct {
	db *gorm.DB
}

func NewShadowRepo(in any) shadow.Repo {
	return &ShadowRepo{db: stores.GetCommonConn(in)}
}

func (p *ShadowRepo) FindByFilter(ctx context.Context, f shadow.Filter) ([]*shadow.Info, error) {
	var results []*DmDeviceShadow
	db := p.fmtFilter(ctx, f).Model(&DmDeviceShadow{})
	err := db.Find(&results).Error
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return ToShadowsDo(results), nil
}

func (p *ShadowRepo) MultiUpdate(ctx context.Context, data []*shadow.Info) error {
	vals := make([]*DmDeviceShadow, len(data))
	for i, d := range data {
		vals[i] = ToShadowPo(d)
	}
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmDeviceShadow{}).Create(vals).Error
	return stores.ErrFmt(err)
}
func (p *ShadowRepo) fmtFilter(ctx context.Context, f shadow.Filter) *gorm.DB {
	db := p.db.WithContext(ctx).Where("product_id = ?", f.ProductID).Where("device_name = ?", f.DeviceName)
	if len(f.DataIDs) != 0 {
		db = db.Where("data_id in ?", f.DataIDs)
	}
	if f.UpdatedDeviceStatus != 0 {
		if f.UpdatedDeviceStatus == shadow.UpdatedDevice {
			db = db.Where("updated_device_time is not null")
		} else {
			db = db.Where("updated_device_time is null")
		}
	}
	return db
}

func (p *ShadowRepo) MultiDelete(ctx context.Context, f shadow.Filter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Unscoped().Delete(&DmDeviceShadow{}).Error
	return errors.IfNotNil(errors.Database, err)
}
