package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将ProtocolScriptDevice全局替换为模型的表名
2. 完善todo
*/

type ProtocolScriptDeviceRepo struct {
	db *gorm.DB
}

func NewProtocolScriptDeviceRepo(in any) *ProtocolScriptDeviceRepo {
	return &ProtocolScriptDeviceRepo{db: stores.GetCommonConn(in)}
}

type ProtocolScriptDeviceFilter struct {
	TriggerSrc int64
	ProductID  string
	DeviceName string
	Status     int64
	WithScript bool
}

func (p ProtocolScriptDeviceRepo) fmtFilter(ctx context.Context, f ProtocolScriptDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TriggerSrc != 0 {
		db = db.Where("trigger_src = ?", f.TriggerSrc)
	}
	if f.ProductID != "" {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name = ?", f.DeviceName)
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.WithScript {
		db = db.Preload("Script")
	}
	return db
}

func (p ProtocolScriptDeviceRepo) Insert(ctx context.Context, data *DmProtocolScriptDevice) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProtocolScriptDeviceRepo) FindOneByFilter(ctx context.Context, f ProtocolScriptDeviceFilter) (*DmProtocolScriptDevice, error) {
	var result DmProtocolScriptDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProtocolScriptDeviceRepo) FindByFilter(ctx context.Context, f ProtocolScriptDeviceFilter, page *stores.PageInfo) ([]*DmProtocolScriptDevice, error) {
	var results []*DmProtocolScriptDevice
	db := p.fmtFilter(ctx, f).Model(&DmProtocolScriptDevice{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProtocolScriptDeviceRepo) CountByFilter(ctx context.Context, f ProtocolScriptDeviceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProtocolScriptDevice{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProtocolScriptDeviceRepo) Update(ctx context.Context, data *DmProtocolScriptDevice) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProtocolScriptDeviceRepo) DeleteByFilter(ctx context.Context, f ProtocolScriptDeviceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProtocolScriptDevice{}).Error
	return stores.ErrFmt(err)
}

func (p ProtocolScriptDeviceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProtocolScriptDevice{}).Error
	return stores.ErrFmt(err)
}
func (p ProtocolScriptDeviceRepo) FindOne(ctx context.Context, id int64) (*DmProtocolScriptDevice, error) {
	var result DmProtocolScriptDevice
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProtocolScriptDeviceRepo) MultiInsert(ctx context.Context, data []*DmProtocolScriptDevice) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProtocolScriptDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d ProtocolScriptDeviceRepo) UpdateWithField(ctx context.Context, f ProtocolScriptDeviceFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProtocolScriptDevice{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
