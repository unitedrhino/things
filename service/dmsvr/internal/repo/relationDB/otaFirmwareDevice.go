package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type OtaFirmwareDeviceRepo struct {
	db *gorm.DB
}

func NewOtaFirmwareDeviceRepo(in any) *OtaFirmwareDeviceRepo {
	return &OtaFirmwareDeviceRepo{db: stores.GetCommonConn(in)}
}

type OtaFirmwareDeviceFilter struct {
	IDs              []int64
	FirmwareID       int64
	JobID            int64
	ProductID        string
	DeviceName       string
	DeviceNames      []string
	WithScheduleTime bool
	//Msg     int64
	Statues      []int64
	SrcVersion   string
	WithFirmware bool
	WithFiles    bool
	IsOnline     int64
}

func (p OtaFirmwareDeviceRepo) fmtFilter(ctx context.Context, f OtaFirmwareDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.FirmwareID != 0 {
		db = db.Where("firmware_id = ?", f.FirmwareID)
	}
	if f.JobID != 0 {
		db = db.Where("job_id = ?", f.JobID)
	}
	if f.IsOnline != 0 && f.ProductID != "" {
		db = db.Where("device_name in (select device_name from dm_device_info where is_online=? and product_id = ?)", f.IsOnline, f.ProductID)
	}
	if f.ProductID != "" {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name like ?", "%"+f.DeviceName+"%")
	}
	if len(f.DeviceNames) != 0 {
		db = db.Where("device_name in ?", f.DeviceNames)
	}
	if f.WithScheduleTime {
		db = db.Where("schedule_time not null")
	}
	if len(f.Statues) != 0 {
		db = db.Where("status in ?", f.Statues)
	}
	if f.SrcVersion != "" {
		db = db.Where("src_version=?", f.SrcVersion)
	}
	if f.WithFirmware {
		db = db.Preload("Firmware")
	}
	if f.WithFiles {
		db = db.Preload("Files")
	}

	return db
}

func (p OtaFirmwareDeviceRepo) Insert(ctx context.Context, data *DmOtaFirmwareDevice) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OtaFirmwareDeviceRepo) FindOneByFilter(ctx context.Context, f OtaFirmwareDeviceFilter) (*DmOtaFirmwareDevice, error) {
	var result DmOtaFirmwareDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaFirmwareDeviceRepo) FindByFilter(ctx context.Context, f OtaFirmwareDeviceFilter, page *def.PageInfo) ([]*DmOtaFirmwareDevice, error) {
	var results []*DmOtaFirmwareDevice
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareDevice{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaFirmwareDeviceRepo) CountByFilter(ctx context.Context, f OtaFirmwareDeviceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareDevice{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OtaFirmwareDeviceRepo) Update(ctx context.Context, data *DmOtaFirmwareDevice) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OtaFirmwareDeviceRepo) DeleteByFilter(ctx context.Context, f OtaFirmwareDeviceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaFirmwareDevice{}).Error
	return stores.ErrFmt(err)
}

func (p OtaFirmwareDeviceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmOtaFirmwareDevice{}).Error
	return stores.ErrFmt(err)
}
func (p OtaFirmwareDeviceRepo) FindOne(ctx context.Context, id int64) (*DmOtaFirmwareDevice, error) {
	var result DmOtaFirmwareDevice
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OtaFirmwareDeviceRepo) MultiInsert(ctx context.Context, data []*DmOtaFirmwareDevice) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaFirmwareDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}

// 批量更新
func (p OtaFirmwareDeviceRepo) BatchUpdateField(ctx context.Context, f OtaFirmwareDeviceFilter, updateData map[string]interface{}) error {
	db := p.fmtFilter(ctx, f)
	err := db.WithContext(ctx).Model(&DmOtaFirmwareDevice{}).Updates(updateData).Error
	return stores.ErrFmt(err)
}
