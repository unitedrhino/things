package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GroupDeviceRepo struct {
	db *gorm.DB
}
type (
	GroupDeviceFilter struct {
		GroupID     int64
		ProductID   string
		DeviceName  string
		WithProduct bool
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
func (m GroupDeviceRepo) MultiDelete(ctx context.Context, groupID int64, data []*DmGroupDevice) error {
	if len(data) < 1 {
		return nil
	}
	scope := func(db *gorm.DB) *gorm.DB {
		for i, d := range data {
			if i == 0 {
				db = db.Where("`productID` = ? and `deviceName` = ?", d.ProductID, d.DeviceName)
				continue
			}
			db = db.Or("`productID` = ? and `deviceName` = ?", d.ProductID, d.DeviceName)
		}
		return db
	}
	db := m.db.WithContext(ctx).Model(&DmGroupDevice{})
	db = db.Where("`groupID`=?", groupID).Where(scope(db))
	err := db.Delete(&DmGroupDevice{}).Error
	return stores.ErrFmt(err)
}

func (p GroupDeviceRepo) fmtFilter(ctx context.Context, f GroupDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	//业务过滤条件
	if f.GroupID != 0 {
		db = db.Where("`groupID`=?", f.GroupID)
	}
	if f.ProductID != "" {
		db = db.Where("`productID`=?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("`deviceName`=?", f.DeviceName)
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
	err := db.Delete(&DmProductInfo{}).Error
	return stores.ErrFmt(err)
}
