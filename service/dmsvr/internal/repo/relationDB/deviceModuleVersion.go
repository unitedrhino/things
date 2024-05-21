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

type DeviceModuleVersionRepo struct {
	db *gorm.DB
}

func NewDeviceModuleVersionRepo(in any) *DeviceModuleVersionRepo {
	return &DeviceModuleVersionRepo{db: stores.GetCommonConn(in)}
}

type DeviceModuleVersionFilter struct {
	ProductID  string
	DeviceName string
	ModuleCode string
}

func (p DeviceModuleVersionRepo) fmtFilter(ctx context.Context, f DeviceModuleVersionFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p DeviceModuleVersionRepo) Insert(ctx context.Context, data *DmDeviceModuleVersion) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeviceModuleVersionRepo) FindOneByFilter(ctx context.Context, f DeviceModuleVersionFilter) (*DmDeviceModuleVersion, error) {
	var result DmDeviceModuleVersion
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DeviceModuleVersionRepo) FindByFilter(ctx context.Context, f DeviceModuleVersionFilter, page *def.PageInfo) ([]*DmDeviceModuleVersion, error) {
	var results []*DmDeviceModuleVersion
	db := p.fmtFilter(ctx, f).Model(&DmDeviceModuleVersion{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeviceModuleVersionRepo) CountByFilter(ctx context.Context, f DeviceModuleVersionFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmDeviceModuleVersion{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeviceModuleVersionRepo) Update(ctx context.Context, data *DmDeviceModuleVersion) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeviceModuleVersionRepo) DeleteByFilter(ctx context.Context, f DeviceModuleVersionFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmDeviceModuleVersion{}).Error
	return stores.ErrFmt(err)
}

func (p DeviceModuleVersionRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmDeviceModuleVersion{}).Error
	return stores.ErrFmt(err)
}
func (p DeviceModuleVersionRepo) FindOne(ctx context.Context, id int64) (*DmDeviceModuleVersion, error) {
	var result DmDeviceModuleVersion
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeviceModuleVersionRepo) MultiInsert(ctx context.Context, data []*DmDeviceModuleVersion) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmDeviceModuleVersion{}).Create(data).Error
	return stores.ErrFmt(err)
}
