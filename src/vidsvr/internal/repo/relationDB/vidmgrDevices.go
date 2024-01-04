package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

// 国标设备接入数据库
type VidmgrDevicesRepo struct {
	db *gorm.DB
}

func NewVidmgrDevicesRepo(in any) *VidmgrDevicesRepo {
	return &VidmgrDevicesRepo{db: stores.GetCommonConn(in)}
}

type VidmgrDevicesFilter struct {
	IDs       []int64
	DeviceIDs []string
	Name      string
	Host      string
}

func (p VidmgrDevicesRepo) fmtFilter(ctx context.Context, f VidmgrDevicesFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if len(f.IDs) != 0 {
		db = db.Where("id in?", f.IDs)
	}
	if len(f.DeviceIDs) != 0 {
		db = db.Where("device_id in?", f.DeviceIDs)
	}
	if f.Name != "" {
		db = db.Where("name =?", f.Name)
	}
	if f.Host != "" {
		db = db.Where("host =?", f.Host)
	}
	return db
}

func (p VidmgrDevicesRepo) Insert(ctx context.Context, data *VidmgrDevices) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p VidmgrDevicesRepo) FindOneByFilter(ctx context.Context, f VidmgrDevicesFilter) (*VidmgrDevices, error) {
	var result VidmgrDevices
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p VidmgrDevicesRepo) Update(ctx context.Context, data *VidmgrDevices) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p VidmgrDevicesRepo) DeleteByFilter(ctx context.Context, f VidmgrDevicesFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&VidmgrDevices{}).Error
	return stores.ErrFmt(err)
}

func (p VidmgrDevicesRepo) FindByFilter(ctx context.Context, f VidmgrDevicesFilter, page *def.PageInfo) ([]*VidmgrDevices, error) {
	var results []*VidmgrDevices
	db := p.fmtFilter(ctx, f).Model(&VidmgrDevices{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrDevicesRepo) CountByFilter(ctx context.Context, f VidmgrDevicesFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&VidmgrDevices{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
