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

type AlarmLogRepo struct {
	db *gorm.DB
}

func NewAlarmLogRepo(in any) *AlarmLogRepo {
	return &AlarmLogRepo{db: stores.GetCommonConn(in)}
}

type AlarmLogFilter struct {
	AlarmRecordID int64 //告警配置ID
	Time          def.TimeRange
}

func (p AlarmLogRepo) fmtFilter(ctx context.Context, f AlarmLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.AlarmRecordID != 0 {
		db = db.Where("alarm_record_id=?", f.AlarmRecordID)
	}
	return db
}

func (p AlarmLogRepo) Insert(ctx context.Context, data *RuleAlarmLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmLogRepo) FindOneByFilter(ctx context.Context, f AlarmLogFilter) (*RuleAlarmLog, error) {
	var result RuleAlarmLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmLogRepo) FindByFilter(ctx context.Context, f AlarmLogFilter, page *stores.PageInfo) ([]*RuleAlarmLog, error) {
	var results []*RuleAlarmLog
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmLogRepo) CountByFilter(ctx context.Context, f AlarmLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmLogRepo) Update(ctx context.Context, data *RuleAlarmLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmLogRepo) DeleteByFilter(ctx context.Context, f AlarmLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&RuleAlarmLog{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&RuleAlarmLog{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmLogRepo) FindOne(ctx context.Context, id int64) (*RuleAlarmLog, error) {
	var result RuleAlarmLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmLogRepo) MultiInsert(ctx context.Context, data []*RuleAlarmLog) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&RuleAlarmLog{}).Create(data).Error
	return stores.ErrFmt(err)
}
