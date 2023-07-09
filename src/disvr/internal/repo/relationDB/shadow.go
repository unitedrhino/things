package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/disvr/internal/domain/shadow"
	"gorm.io/gorm"
)

type ShadowRepo struct {
	db *gorm.DB
}

func NewShadowRepo(in any) shadow.Repo {
	return &ShadowRepo{db: stores.GetTenantConn(in)}
}

func (p *ShadowRepo) FindByFilter(ctx context.Context, f shadow.Filter) ([]*shadow.Info, error) {
	var results []*DiDeviceShadow
	db := p.fmtFilter(ctx, f).Model(&DiDeviceShadow{})
	err := db.Find(&results).Error
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return ToShadowsDo(results), nil
}

func (p *ShadowRepo) MultiUpdate(ctx context.Context, data []*shadow.Info) error {
	vals := make([]*DiDeviceShadow, len(data))
	for i, d := range data {
		vals[i] = ToShadowPo(d)
	}
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range vals {
			err := tx.Unscoped().Delete(&DiDeviceShadow{}, "productID = ? and deviceName = ? and dataID = ?",
				v.ProductID, v.DeviceName, v.DataID).Error
			if err != nil {
				return errors.Database.AddDetail(err)
			}
			err = tx.Save(v).Error
			if err != nil {
				return errors.Database.AddDetail(err)
			}
		}
		return nil
	})
	return err
}
func (p *ShadowRepo) fmtFilter(ctx context.Context, f shadow.Filter) *gorm.DB {
	db := p.db.WithContext(ctx).Where("productID = ?", f.ProductID).Where("deviceName = ?", f.DeviceName)
	if len(f.DataIDs) != 0 {
		db = db.Where("dataID in ?", f.DataIDs)
	}
	if f.UpdatedDeviceStatus != 0 {
		if f.UpdatedDeviceStatus == shadow.UpdatedDevice {
			db = db.Where("updatedDeviceTime is not null")
		} else {
			db = db.Where("updatedDeviceTime is null")
		}
	}
	return db
}

func (p *ShadowRepo) MultiDelete(ctx context.Context, f shadow.Filter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Unscoped().Delete(&DiDeviceShadow{}).Error
	return errors.IfNotNil(errors.Database, err)
}
