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

type SceneIfTriggerRepo struct {
	db *gorm.DB
}

func NewSceneIfTriggerRepo(in any) *SceneIfTriggerRepo {
	return &SceneIfTriggerRepo{db: stores.GetCommonConn(in)}
}

type SceneIfTriggerFilter struct {
	SceneID     int64
	Status      int64
	ExecAt      *stores.Cmp
	LastRunTime *stores.Cmp
	Repeat      *stores.Cmp
	Type        string
}

func (p SceneIfTriggerRepo) fmtFilter(ctx context.Context, f SceneIfTriggerFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.ExecAt.Where(db, "timer_exec_at")
	db = f.LastRunTime.Where(db, "last_run_time")
	db = f.Repeat.Where(db, "timer_exec_repeat")
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.Type != "" {
		db = db.Where("type = ?", f.Type)
	}
	if f.SceneID != 0 {
		db = db.Where("scene_id = ?", f.SceneID)
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
	err := db.Preload("SceneInfo").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneIfTriggerRepo) FindByFilter(ctx context.Context, f SceneIfTriggerFilter, page *def.PageInfo) ([]*UdSceneIfTrigger, error) {
	var results []*UdSceneIfTrigger
	db := p.fmtFilter(ctx, f).Model(&UdSceneIfTrigger{})
	db = page.ToGorm(db)
	err := db.Preload("SceneInfo").Find(&results).Error
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
	err := p.db.WithContext(ctx).Preload("SceneInfo").Where("id = ?", id).First(&result).Error
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
