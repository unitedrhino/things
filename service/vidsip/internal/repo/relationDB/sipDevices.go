package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
)

// 国标设备接入数据库
type SipDevicesRepo struct {
	db *gorm.DB
}

func NewSipDevicesRepo(in any) *SipDevicesRepo {
	return &SipDevicesRepo{db: stores.GetCommonConn(in)}
}

type SipDevicesFilter struct {
	DeviceIDs []string
	DeviceID  string
	Name      string
	Host      string
}

func (p SipDevicesRepo) fmtFilter(ctx context.Context, f SipDevicesFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if len(f.DeviceIDs) != 0 {
		db = db.Where("device_id in?", f.DeviceIDs)
	}

	if f.DeviceID != "" {
		db = db.Where("device_id  =?", f.DeviceID)
	}
	if f.Name != "" {
		db = db.Where("name =?", f.Name)
	}
	if f.Host != "" {
		db = db.Where("host =?", f.Host)
	}
	return db
}

func (p SipDevicesRepo) Insert(ctx context.Context, data *SipDevices) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SipDevicesRepo) FindOneByFilter(ctx context.Context, f SipDevicesFilter) (*SipDevices, error) {
	var result SipDevices
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p SipDevicesRepo) Update(ctx context.Context, data *SipDevices) error {
	err := p.db.WithContext(ctx).Where("device_id = ?", data.DeviceID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SipDevicesRepo) DeleteByFilter(ctx context.Context, f SipDevicesFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SipDevices{}).Error
	return stores.ErrFmt(err)
}

func (p SipDevicesRepo) FindByFilter(ctx context.Context, f SipDevicesFilter, page *def.PageInfo) ([]*SipDevices, error) {
	var results []*SipDevices
	db := p.fmtFilter(ctx, f).Model(&SipDevices{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SipDevicesRepo) CountByFilter(ctx context.Context, f SipDevicesFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SipDevices{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
