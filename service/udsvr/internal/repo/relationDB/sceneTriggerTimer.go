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

type SceneTriggerTimerRepo struct {
	db *gorm.DB
}

func NewSceneTriggerTimerRepo(in any) *SceneTriggerTimerRepo {
	return &SceneTriggerTimerRepo{db: stores.GetCommonConn(in)}
}

type SceneTriggerTimerFilter struct {
	SceneID     int64
	Status      int64
	ExecAt      *stores.Cmp
	LastRunTime *stores.Cmp
	Repeat      *stores.Cmp
}

func (p SceneTriggerTimerRepo) fmtFilter(ctx context.Context, f SceneTriggerTimerFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.ExecAt.Where(db, "exec_at")
	db = f.LastRunTime.Where(db, "last_run_time")
	db = f.Repeat.Where(db, "exec_repeat")
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.SceneID != 0 {
		db = db.Where("scene_id = ?", f.SceneID)
	}
	return db
}

func (p SceneTriggerTimerRepo) Insert(ctx context.Context, data *UdSceneTriggerTimer) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SceneTriggerTimerRepo) FindOneByFilter(ctx context.Context, f SceneTriggerTimerFilter) (*UdSceneTriggerTimer, error) {
	var result UdSceneTriggerTimer
	db := p.fmtFilter(ctx, f)
	err := db.Preload("SceneInfo").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneTriggerTimerRepo) FindByFilter(ctx context.Context, f SceneTriggerTimerFilter, page *def.PageInfo) ([]*UdSceneTriggerTimer, error) {
	var results []*UdSceneTriggerTimer
	db := p.fmtFilter(ctx, f).Model(&UdSceneTriggerTimer{})
	db = page.ToGorm(db)
	err := db.Preload("SceneInfo").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SceneTriggerTimerRepo) CountByFilter(ctx context.Context, f SceneTriggerTimerFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdSceneTriggerTimer{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneTriggerTimerRepo) Update(ctx context.Context, data *UdSceneTriggerTimer) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SceneTriggerTimerRepo) DeleteByFilter(ctx context.Context, f SceneTriggerTimerFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdSceneTriggerTimer{}).Error
	return stores.ErrFmt(err)
}

func (p SceneTriggerTimerRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdSceneTriggerTimer{}).Error
	return stores.ErrFmt(err)
}
func (p SceneTriggerTimerRepo) FindOne(ctx context.Context, id int64) (*UdSceneTriggerTimer, error) {
	var result UdSceneTriggerTimer
	err := p.db.WithContext(ctx).Preload("SceneInfo").Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneTriggerTimerRepo) MultiInsert(ctx context.Context, data []*UdSceneTriggerTimer) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdSceneTriggerTimer{}).Create(data).Error
	return stores.ErrFmt(err)
}
