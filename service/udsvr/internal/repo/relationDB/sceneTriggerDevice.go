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

type SceneTriggerDeviceRepo struct {
	db *gorm.DB
}

func NewSceneTriggerDeviceRepo(in any) *SceneTriggerDeviceRepo {
	return &SceneTriggerDeviceRepo{db: stores.GetCommonConn(in)}
}

type SceneTriggerDeviceFilter struct {
	SceneID int64
}

func (p SceneTriggerDeviceRepo) fmtFilter(ctx context.Context, f SceneTriggerDeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.SceneID != 0 {
		db = db.Where("scene_id = ?", f.SceneID)
	}
	return db
}

func (p SceneTriggerDeviceRepo) Insert(ctx context.Context, data *UdSceneTriggerDevice) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SceneTriggerDeviceRepo) FindOneByFilter(ctx context.Context, f SceneTriggerDeviceFilter) (*UdSceneTriggerDevice, error) {
	var result UdSceneTriggerDevice
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneTriggerDeviceRepo) FindByFilter(ctx context.Context, f SceneTriggerDeviceFilter, page *def.PageInfo) ([]*UdSceneTriggerDevice, error) {
	var results []*UdSceneTriggerDevice
	db := p.fmtFilter(ctx, f).Model(&UdSceneTriggerDevice{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SceneTriggerDeviceRepo) CountByFilter(ctx context.Context, f SceneTriggerDeviceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdSceneTriggerDevice{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneTriggerDeviceRepo) Update(ctx context.Context, data *UdSceneTriggerDevice) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SceneTriggerDeviceRepo) DeleteByFilter(ctx context.Context, f SceneTriggerDeviceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdSceneTriggerDevice{}).Error
	return stores.ErrFmt(err)
}

func (p SceneTriggerDeviceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdSceneTriggerDevice{}).Error
	return stores.ErrFmt(err)
}
func (p SceneTriggerDeviceRepo) FindOne(ctx context.Context, id int64) (*UdSceneTriggerDevice, error) {
	var result UdSceneTriggerDevice
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneTriggerDeviceRepo) MultiInsert(ctx context.Context, data []*UdSceneTriggerDevice) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdSceneTriggerDevice{}).Create(data).Error
	return stores.ErrFmt(err)
}
