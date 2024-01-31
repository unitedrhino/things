package relationDB

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type AlarmDealRecordRepo struct {
	db *gorm.DB
}

func NewAlarmDealRecordRepo(in any) *AlarmDealRecordRepo {
	return &AlarmDealRecordRepo{db: stores.GetCommonConn(in)}
}

type AlarmDealRecordFilter struct {
	AlarmRecordID int64 //告警配置ID
	Time          def.TimeRange
}

func (p AlarmDealRecordRepo) fmtFilter(ctx context.Context, f AlarmDealRecordFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.AlarmRecordID != 0 {
		db = db.Where("alarm_record_id=?", f.AlarmRecordID)
	}
	return db
}

func (p AlarmDealRecordRepo) Insert(ctx context.Context, data *RuleAlarmDealRecord) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmDealRecordRepo) FindOneByFilter(ctx context.Context, f AlarmDealRecordFilter) (*RuleAlarmDealRecord, error) {
	var result RuleAlarmDealRecord
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmDealRecordRepo) FindByFilter(ctx context.Context, f AlarmDealRecordFilter, page *def.PageInfo) ([]*RuleAlarmDealRecord, error) {
	var results []*RuleAlarmDealRecord
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmDealRecord{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmDealRecordRepo) CountByFilter(ctx context.Context, f AlarmDealRecordFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmDealRecord{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmDealRecordRepo) Update(ctx context.Context, data *RuleAlarmDealRecord) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmDealRecordRepo) DeleteByFilter(ctx context.Context, f AlarmDealRecordFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&RuleAlarmDealRecord{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmDealRecordRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&RuleAlarmDealRecord{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmDealRecordRepo) FindOne(ctx context.Context, id int64) (*RuleAlarmDealRecord, error) {
	var result RuleAlarmDealRecord
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmDealRecordRepo) MultiInsert(ctx context.Context, data []*RuleAlarmDealRecord) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&RuleAlarmDealRecord{}).Create(data).Error
	return stores.ErrFmt(err)
}
