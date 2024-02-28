package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GroupDeviceRepo struct {
	db *gorm.DB
}
type (
	GroupDeviceFilter struct {
		GroupIDs    []int64
		ProductID   string
		DeviceName  string
		WithProduct bool
		WithDevice  bool
	}
)

func NewGroupDeviceRepo(in any) *GroupDeviceRepo {
	return &GroupDeviceRepo{db: stores.GetCommonConn(in)}
}

// 批量插入 LightStrategyDevice 记录
func (m GroupDeviceRepo) MultiInsert(ctx context.Context, data []*DmGroupDevice) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmGroupDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (m GroupDeviceRepo) MultiDelete(ctx context.Context, groupID int64, data []*devices.Core) error {
	if len(data) < 1 {
		return nil
	}
	scope := func(db *gorm.DB) *gorm.DB {
		for i, d := range data {
			if i == 0 {
				db = db.Where("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
				continue
			}
			db = db.Or("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
		}
		return db
	}
	db := m.db.WithContext(ctx).Model(&DmGroupDevice{})
	db = db.Where("group_id=?", groupID).Where(scope(db))
	err := db.Delete(&DmGroupDevice{}).Error
	return stores.ErrFmt(err)
}

func (p GroupDeviceRepo) MultiUpdate(ctx context.Context, groupID int64, devices []*DmGroupDevice) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewGroupDeviceRepo(tx)
		err := rm.DeleteByFilter(ctx, GroupDeviceFilter{GroupIDs: []int64{groupID}})
		if err != nil {
			return err
		}
		if len(devices) != 0 {
			err = rm.MultiInsert(ctx, devices)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return stores.ErrFmt(err)
}

func (p GroupDeviceRepo) fmtFilter(ctx context.Context, f GroupDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	if f.WithDevice {
		db = db.Preload("Device")
	}
	//业务过滤条件
	if len(f.GroupIDs) != 0 {
		db = db.Where("group_id in ?", f.GroupIDs)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name=?", f.DeviceName)
	}
	return db
}

func (p GroupDeviceRepo) FindByFilter(ctx context.Context, f GroupDeviceFilter, page *def.PageInfo) ([]*DmGroupDevice, error) {
	var results []*DmGroupDevice
	db := p.fmtFilter(ctx, f).Model(&DmGroupDevice{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p GroupDeviceRepo) FindOneByFilter(ctx context.Context, f GroupDeviceFilter) (*DmGroupDevice, error) {
	var result DmGroupDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p GroupDeviceRepo) CountByFilter(ctx context.Context, f GroupDeviceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmGroupDevice{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
func (p GroupDeviceRepo) DeleteByFilter(ctx context.Context, f GroupDeviceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmGroupDevice{}).Error
	return stores.ErrFmt(err)
}
