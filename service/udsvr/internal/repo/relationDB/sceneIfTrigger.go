package relationDB

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type SceneIfTriggerRepo struct {
	db *gorm.DB
}

func NewSceneIfTriggerRepo(in any) *SceneIfTriggerRepo {
	return &SceneIfTriggerRepo{db: stores.GetCommonConn(in)}
}

type SceneIfTriggerFilter struct {
	ID               int64
	SceneID          int64
	Status           int64
	ExecAt           *stores.Cmp
	LastRunTime      *stores.Cmp
	ExecRepeat       *stores.Cmp
	RepeatType       scene.RepeatType
	ExecLoopStart    *stores.Cmp
	ExecLoopEnd      *stores.Cmp
	ExecType         *stores.Cmp
	Type             string
	ProjectID        *stores.Cmp
	AreaID           *stores.Cmp
	Device           *devices.Core
	DataID           string
	FirstTriggerTime *stores.Cmp
	StateKeepType    scene.StateKeepType
	StateKeepValue   *stores.Cmp
}

func (p SceneIfTriggerRepo) fmtFilter(ctx context.Context, f SceneIfTriggerFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.ProjectID.Where(db, "project_id")
	db = f.AreaID.Where(db, "area_id")
	db = f.ExecAt.Where(db, "timer_exec_at")
	db = f.LastRunTime.Where(db, "last_run_time")
	db = f.ExecRepeat.Where(db, "timer_exec_repeat")
	db = f.ExecLoopStart.Where(db, "timer_exec_loop_start")
	db = f.ExecLoopEnd.Where(db, "timer_exec_loop_end")
	db = f.FirstTriggerTime.Where(db, "device_first_trigger_time")
	db = f.StateKeepValue.Where(db, "device_state_keep_value")
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	db = f.ExecType.Where(db, "timer_exec_type")
	if f.RepeatType != "" {
		db = db.Where("timer_repeat_type = ?", f.RepeatType)
	}
	if f.Type != "" {
		db = db.Where("type = ?", f.Type)
	}
	if f.StateKeepType != "" {
		db = db.Where("device_state_keep_type = ?", f.StateKeepType)
	}
	if f.SceneID != 0 {
		db = db.Where("scene_id = ?", f.SceneID)
	}
	if f.Device != nil {
		db = db.Where("device_select_type=? and device_product_id=? and device_device_name=? and device_data_id=?",
			scene.SelectDeviceFixed, f.Device.ProductID, f.Device.DeviceName, f.DataID).
			Or("device_select_type=? and device_product_id=? and device_data_id=?", scene.SelectorDeviceAll,
				f.Device.ProductID, f.DataID)
	}
	return db
}

func (p SceneIfTriggerRepo) Insert(ctx context.Context, data *UdSceneIfTrigger) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SceneIfTriggerRepo) FindOneByFilter(ctx context.Context, f SceneIfTriggerFilter) (*UdSceneIfTrigger, error) {
	var result UdSceneIfTrigger
	db := p.fmtFilter(ctx, f)
	err := db.Preload("SceneInfo").Preload("SceneInfo.Actions").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneIfTriggerRepo) FindByFilter(ctx context.Context, f SceneIfTriggerFilter, page *stores.PageInfo) ([]*UdSceneIfTrigger, error) {
	var results []*UdSceneIfTrigger
	db := p.fmtFilter(ctx, f).Model(&UdSceneIfTrigger{})
	db = page.ToGorm(db)
	err := db.Preload("SceneInfo").Preload("SceneInfo.Actions").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SceneIfTriggerRepo) CountByFilter(ctx context.Context, f SceneIfTriggerFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdSceneIfTrigger{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneIfTriggerRepo) Update(ctx context.Context, data *UdSceneIfTrigger) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d SceneIfTriggerRepo) UpdateWithField(ctx context.Context, f SceneIfTriggerFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UdSceneIfTrigger{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p SceneIfTriggerRepo) DeleteByFilter(ctx context.Context, f SceneIfTriggerFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdSceneIfTrigger{}).Error
	return stores.ErrFmt(err)
}

func (p SceneIfTriggerRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdSceneIfTrigger{}).Error
	return stores.ErrFmt(err)
}
func (p SceneIfTriggerRepo) FindOne(ctx context.Context, id int64) (*UdSceneIfTrigger, error) {
	var result UdSceneIfTrigger
	err := p.db.WithContext(ctx).Preload("SceneInfo").Preload("SceneInfo.Actions").Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneIfTriggerRepo) MultiInsert(ctx context.Context, data []*UdSceneIfTrigger) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdSceneIfTrigger{}).Create(data).Error
	return stores.ErrFmt(err)
}
