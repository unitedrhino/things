package relationDB

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将DeviceProfile全局替换为模型的表名
2. 完善todo
*/

type DeviceProfileRepo struct {
	db *gorm.DB
}

func NewDeviceProfileRepo(in any) *DeviceProfileRepo {
	return &DeviceProfileRepo{db: stores.GetCommonConn(in)}
}

type DeviceProfileFilter struct {
	Codes  []string
	Code   string
	Device devices.Core
}

func (p DeviceProfileRepo) fmtFilter(ctx context.Context, f DeviceProfileFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.Codes) != 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	db = db.Where("product_id =? and device_name=?",
		f.Device.ProductID, f.Device.DeviceName)
	return db
}

func (p DeviceProfileRepo) Insert(ctx context.Context, data *DmDeviceProfile) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeviceProfileRepo) FindOneByFilter(ctx context.Context, f DeviceProfileFilter) (*DmDeviceProfile, error) {
	var result DmDeviceProfile
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DeviceProfileRepo) FindByFilter(ctx context.Context, f DeviceProfileFilter, page *stores.PageInfo) ([]*DmDeviceProfile, error) {
	var results []*DmDeviceProfile
	db := p.fmtFilter(ctx, f).Model(&DmDeviceProfile{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeviceProfileRepo) CountByFilter(ctx context.Context, f DeviceProfileFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmDeviceProfile{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeviceProfileRepo) Update(ctx context.Context, data *DmDeviceProfile) error {
	err := p.db.WithContext(ctx).Where("product_id = ? and device_name=? and code = ?",
		data.ProductID, data.DeviceName, data.Code).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeviceProfileRepo) DeleteByFilter(ctx context.Context, f DeviceProfileFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmDeviceProfile{}).Error
	return stores.ErrFmt(err)
}

func (p DeviceProfileRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmDeviceProfile{}).Error
	return stores.ErrFmt(err)
}
func (p DeviceProfileRepo) FindOne(ctx context.Context, id int64) (*DmDeviceProfile, error) {
	var result DmDeviceProfile
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeviceProfileRepo) MultiInsert(ctx context.Context, data []*DmDeviceProfile) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmDeviceProfile{}).Create(data).Error
	return stores.ErrFmt(err)
}
