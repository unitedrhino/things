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

type AlarmRecordRepo struct {
	db *gorm.DB
}

func NewAlarmRecordRepo(in any) *AlarmRecordRepo {
	return &AlarmRecordRepo{db: stores.GetCommonConn(in)}
}

type AlarmRecordFilter struct {
	AlarmID     int64 // 告警配置ID
	TriggerType int64
	ProductID   string
	DeviceName  string
	Time        def.TimeRange
}

func (p AlarmRecordRepo) fmtFilter(ctx context.Context, f AlarmRecordFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.AlarmID != 0 {
		db = db.Where("alarm_id=?", f.AlarmID)
	}
	if f.TriggerType != 0 {
		db = db.Where("trigger_type=?", f.TriggerType)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name=?", f.DeviceName)
	}
	return db
}

func (p AlarmRecordRepo) Insert(ctx context.Context, data *RuleAlarmRecord) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmRecordRepo) FindOneByFilter(ctx context.Context, f AlarmRecordFilter) (*RuleAlarmRecord, error) {
	var result RuleAlarmRecord
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmRecordRepo) FindByFilter(ctx context.Context, f AlarmRecordFilter, page *stores.PageInfo) ([]*RuleAlarmRecord, error) {
	var results []*RuleAlarmRecord
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmRecord{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmRecordRepo) CountByFilter(ctx context.Context, f AlarmRecordFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&RuleAlarmRecord{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmRecordRepo) Update(ctx context.Context, data *RuleAlarmRecord) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmRecordRepo) DeleteByFilter(ctx context.Context, f AlarmRecordFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&RuleAlarmRecord{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmRecordRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&RuleAlarmRecord{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmRecordRepo) FindOne(ctx context.Context, id int64) (*RuleAlarmRecord, error) {
	var result RuleAlarmRecord
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmRecordRepo) MultiInsert(ctx context.Context, data []*RuleAlarmRecord) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&RuleAlarmRecord{}).Create(data).Error
	return stores.ErrFmt(err)
}
