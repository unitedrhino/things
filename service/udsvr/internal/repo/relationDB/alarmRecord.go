package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
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

type AlarmRecordRepo struct {
	db *gorm.DB
}

func NewAlarmRecordRepo(in any) *AlarmRecordRepo {
	return &AlarmRecordRepo{db: stores.GetCommonConn(in)}
}

type AlarmRecordFilter struct {
	AlarmID      int64 // 告警配置ID
	TriggerType  scene.TriggerType
	ProductID    string
	DeviceName   string
	DealStatus   scene.AlarmDealStatus
	DealStatuses []scene.AlarmDealStatus
	Time         def.TimeRange
}

func (p AlarmRecordRepo) fmtFilter(ctx context.Context, f AlarmRecordFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.AlarmID != 0 {
		db = db.Where("alarm_id=?", f.AlarmID)
	}
	if f.DealStatus != 0 {
		db = db.Where("deal_status=?", f.DealStatus)
	}
	if len(f.DealStatuses) != 0 {
		db = db.Where("deal_status in ?", f.DealStatuses)
	}
	if f.TriggerType != "" {
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

func (p AlarmRecordRepo) Insert(ctx context.Context, data *UdAlarmRecord) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmRecordRepo) FindOneByFilter(ctx context.Context, f AlarmRecordFilter) (*UdAlarmRecord, error) {
	var result UdAlarmRecord
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmRecordRepo) FindByFilter(ctx context.Context, f AlarmRecordFilter, page *def.PageInfo) ([]*UdAlarmRecord, error) {
	var results []*UdAlarmRecord
	db := p.fmtFilter(ctx, f).Model(&UdAlarmRecord{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmRecordRepo) CountByFilter(ctx context.Context, f AlarmRecordFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdAlarmRecord{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmRecordRepo) Update(ctx context.Context, data *UdAlarmRecord) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmRecordRepo) DeleteByFilter(ctx context.Context, f AlarmRecordFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdAlarmRecord{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmRecordRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdAlarmRecord{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmRecordRepo) FindOne(ctx context.Context, id int64) (*UdAlarmRecord, error) {
	var result UdAlarmRecord
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmRecordRepo) MultiInsert(ctx context.Context, data []*UdAlarmRecord) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdAlarmRecord{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d AlarmRecordRepo) UpdateWithField(ctx context.Context, f AlarmRecordFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UdAlarmRecord{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
