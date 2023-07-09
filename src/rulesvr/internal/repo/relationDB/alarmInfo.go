package relationDB

import (
	"context"
	"fmt"
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

type AlarmInfoRepo struct {
	db *gorm.DB
}

func NewAlarmInfoRepo(in any) *AlarmInfoRepo {
	return &AlarmInfoRepo{db: stores.GetCommonConn(in)}
}

type AlarmInfoFilter struct {
	Name     string //名字
	SceneID  int64  // 场景ID
	AlarmIDs []int64
}

func (p AlarmInfoRepo) fmtFilter(ctx context.Context, f AlarmInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name=?", f.Name)
	}
	if len(f.AlarmIDs) != 0 {
		db = db.Where(fmt.Sprintf("id in (%v)", stores.ArrayToSql(f.AlarmIDs)))
	}
	if f.SceneID != 0 {
		table := RuleAlarmInfo{}
		db = db.Joins(fmt.Sprintf("left join `rule_alarm_scene` as ras on ras.alarmID=%s.id", table.TableName()))
		db = db.Where("ras.sceneID=?", f.SceneID)
	}
	return db
}

func (p AlarmInfoRepo) Insert(ctx context.Context, data *RuleAlarmInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmInfoRepo) FindOneByFilter(ctx context.Context, f AlarmInfoFilter) (*RuleAlarmInfo, error) {
	var result RuleAlarmInfo
	db := p.fmtFilter(ctx, f)
	table := RuleAlarmInfo{}
	err := db.Select(table.TableName() + ".*").First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmInfoRepo) FindByFilter(ctx context.Context, f AlarmInfoFilter, page *def.PageInfo) ([]*RuleAlarmInfo, error) {
	var results []*RuleAlarmInfo
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmInfo{})
	db = page.ToGorm(db)
	table := RuleAlarmInfo{}
	err := db.Select(table.TableName() + ".*").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmInfoRepo) CountByFilter(ctx context.Context, f AlarmInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmInfoRepo) Update(ctx context.Context, data *RuleAlarmInfo) error {
	err := p.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmInfoRepo) DeleteByFilter(ctx context.Context, f AlarmInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&RuleAlarmInfo{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("`id` = ?", id).Delete(&RuleAlarmInfo{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmInfoRepo) FindOne(ctx context.Context, id int64) (*RuleAlarmInfo, error) {
	var result RuleAlarmInfo
	err := p.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmInfoRepo) MultiInsert(ctx context.Context, data []*RuleAlarmInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&RuleAlarmInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
