package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type AlarmSceneRepo struct {
	db *gorm.DB
}

func NewAlarmSceneRepo(in any) *AlarmSceneRepo {
	return &AlarmSceneRepo{db: stores.GetCommonConn(in)}
}

type AlarmSceneFilter struct {
	AlarmID int64 // 告警配置ID
	SceneID int64 // 场景ID
}

func (p AlarmSceneRepo) fmtFilter(ctx context.Context, f AlarmSceneFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.SceneID != 0 {
		db = db.Where("scene_id=?", f.SceneID)
	}
	if f.AlarmID != 0 {
		db = db.Where("alarm_id=?", f.AlarmID)
	}
	return db
}

func (p AlarmSceneRepo) Insert(ctx context.Context, data *RuleAlarmScene) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmSceneRepo) FindOneByFilter(ctx context.Context, f AlarmSceneFilter) (*RuleAlarmScene, error) {
	var result RuleAlarmScene
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmSceneRepo) FindByFilter(ctx context.Context, f AlarmSceneFilter, page *def.PageInfo) ([]*RuleAlarmScene, error) {
	var results []*RuleAlarmScene
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmScene{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmSceneRepo) CountByFilter(ctx context.Context, f AlarmSceneFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmScene{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmSceneRepo) Update(ctx context.Context, data *RuleAlarmScene) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmSceneRepo) DeleteByFilter(ctx context.Context, f AlarmSceneFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&RuleAlarmScene{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmSceneRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&RuleAlarmScene{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmSceneRepo) FindOne(ctx context.Context, id int64) (*RuleAlarmScene, error) {
	var result RuleAlarmScene
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmSceneRepo) MultiInsert(ctx context.Context, alarmID int64, sceneIDs []int64) error {
	var datas []*RuleAlarmScene
	for _, v := range sceneIDs {
		datas = append(datas, &RuleAlarmScene{
			AlarmID: alarmID,
			SceneID: v,
		})
	}
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&RuleAlarmScene{}).Create(datas).Error
	return stores.ErrFmt(err)
}
