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

type AlarmInfoRepo struct {
	db *gorm.DB
}

func NewAlarmInfoRepo(in any) *AlarmInfoRepo {
	return &AlarmInfoRepo{db: stores.GetCommonConn(in)}
}

type AlarmInfoFilter struct {
	Name string
}

func (p AlarmInfoRepo) fmtFilter(ctx context.Context, f AlarmInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	return db
}

func (p AlarmInfoRepo) Insert(ctx context.Context, data *UdAlarmInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AlarmInfoRepo) FindOneByFilter(ctx context.Context, f AlarmInfoFilter) (*UdAlarmInfo, error) {
	var result UdAlarmInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AlarmInfoRepo) FindByFilter(ctx context.Context, f AlarmInfoFilter, page *def.PageInfo) ([]*UdAlarmInfo, error) {
	var results []*UdAlarmInfo
	db := p.fmtFilter(ctx, f).Model(&UdAlarmInfo{})
	db = page.ToGorm(db)
	err := db.Preload("Scenes").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AlarmInfoRepo) CountByFilter(ctx context.Context, f AlarmInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdAlarmInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AlarmInfoRepo) Update(ctx context.Context, data *UdAlarmInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AlarmInfoRepo) DeleteByFilter(ctx context.Context, f AlarmInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdAlarmInfo{}).Error
	return stores.ErrFmt(err)
}

func (p AlarmInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdAlarmInfo{}).Error
	return stores.ErrFmt(err)
}
func (p AlarmInfoRepo) FindOne(ctx context.Context, id int64) (*UdAlarmInfo, error) {
	var result UdAlarmInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AlarmInfoRepo) MultiInsert(ctx context.Context, data []*UdAlarmInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdAlarmInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
