package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type SceneActionRepo struct {
	db *gorm.DB
}

func NewSceneActionRepo(in any) *SceneActionRepo {
	return &SceneActionRepo{db: stores.GetCommonConn(in)}
}

type SceneActionFilter struct {
	SceneID          int64
	Status           scene.Status
	Statuses         []scene.Status
	ProductID        string
	DeviceName       string
	DeviceAreaID     *stores.Cmp
	ActionDeviceType scene.ActionDeviceType
	DeviceSelectType scene.SelectType
}

func (p SceneActionRepo) fmtFilter(ctx context.Context, f SceneActionFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.DeviceAreaID.Where(db, "device_area_id")
	if f.SceneID != 0 {
		db = db.Where("scene_id = ?", f.SceneID)
	}
	if f.ProductID != "" {
		db = db.Where("device_product_id = ?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("device_device_name = ?", f.DeviceName)
	}
	if f.ActionDeviceType != "" {
		db = db.Where("device_action_device_type = ?", f.ActionDeviceType)
	}
	if f.DeviceSelectType != "" {
		db = db.Where("device_select_type = ?", f.DeviceSelectType)
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if len(f.Statuses) != 0 {
		db = db.Where("status in ?", f.Statuses)
	}
	return db
}

func (p SceneActionRepo) Insert(ctx context.Context, data *UdSceneThenAction) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SceneActionRepo) FindOneByFilter(ctx context.Context, f SceneActionFilter) (*UdSceneThenAction, error) {
	var result UdSceneThenAction
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneActionRepo) FindByFilter(ctx context.Context, f SceneActionFilter, page *stores.PageInfo) (
	[]*UdSceneThenAction, error) {
	var results []*UdSceneThenAction
	db := p.fmtFilter(ctx, f).Model(&UdSceneThenAction{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SceneActionRepo) CountByFilter(ctx context.Context, f SceneActionFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdSceneThenAction{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneActionRepo) Update(ctx context.Context, data *UdSceneThenAction) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d SceneActionRepo) UpdateWithField(ctx context.Context, f SceneActionFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UdSceneThenAction{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p SceneActionRepo) DeleteByFilter(ctx context.Context, f SceneActionFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdSceneThenAction{}).Error
	return stores.ErrFmt(err)
}

func (p SceneActionRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdSceneThenAction{}).Error
	return stores.ErrFmt(err)
}
func (p SceneActionRepo) FindOne(ctx context.Context, id int64) (*UdSceneThenAction, error) {
	var result UdSceneThenAction
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneActionRepo) MultiInsert(ctx context.Context, data []*UdSceneThenAction) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdSceneThenAction{}).Create(data).Error
	return stores.ErrFmt(err)
}
