package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/stores"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type SceneInfoRepo struct {
	db *gorm.DB
}

func NewSceneInfoRepo(in any) *SceneInfoRepo {
	return &SceneInfoRepo{db: stores.GetCommonConn(in)}
}

func (p SceneInfoRepo) fmtFilter(ctx context.Context, f scene.InfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.TriggerType != "" {
		db = db.Where("trigger_type = ?", f.TriggerType)
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.AlarmID != 0 {
		table := RuleSceneInfo{}
		db = db.Joins(fmt.Sprintf("left join `rule_alarm_scene` as ras on ras.sceneID=%s.id", table.TableName()))
		db = db.Where("ras.alarm_id=?", f.AlarmID)
	}
	return db
}

func (p SceneInfoRepo) Insert(ctx context.Context, data *scene.Info) (id int64, err error) {
	po := SceneInfoDoToPo(data)
	result := p.db.WithContext(ctx).Create(po)
	return po.ID, stores.ErrFmt(result.Error)
}

func (p SceneInfoRepo) FindOneByFilter(ctx context.Context, f scene.InfoFilter) (*scene.Info, error) {
	var result RuleSceneInfo
	db := p.fmtFilter(ctx, f)
	table := RuleSceneInfo{}
	err := db.Select(table.TableName() + ".*").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return SceneInfoPoToDo(&result), nil
}
func (p SceneInfoRepo) FindByFilter(ctx context.Context, f scene.InfoFilter, page *def.PageInfo) (scene.Infos, error) {
	var results []*RuleSceneInfo
	db := p.fmtFilter(ctx, f).Model(&RuleSceneInfo{})
	db = page.ToGorm(db)
	table := RuleSceneInfo{}
	err := db.Select(table.TableName() + ".*").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return SceneInfoPoToDos(results), nil
}

func (p SceneInfoRepo) CountByFilter(ctx context.Context, f scene.InfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&RuleSceneInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneInfoRepo) Update(ctx context.Context, data *scene.Info) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(SceneInfoDoToPo(data)).Error
	return stores.ErrFmt(err)
}

func (p SceneInfoRepo) DeleteByFilter(ctx context.Context, f scene.InfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&RuleSceneInfo{}).Error
	return stores.ErrFmt(err)
}

func (p SceneInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&RuleSceneInfo{}).Error
	return stores.ErrFmt(err)
}
func (p SceneInfoRepo) FindOne(ctx context.Context, id int64) (*scene.Info, error) {
	var result RuleSceneInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return SceneInfoPoToDo(&result), nil
}

func (p SceneInfoRepo) FindOneByName(ctx context.Context, name string) (*scene.Info, error) {
	var result RuleSceneInfo
	err := p.db.WithContext(ctx).Where("name = ?", name).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return SceneInfoPoToDo(&result), nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneInfoRepo) MultiInsert(ctx context.Context, data []*RuleSceneInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&RuleSceneInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
