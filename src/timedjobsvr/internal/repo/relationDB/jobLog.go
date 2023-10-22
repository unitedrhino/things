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

type JobLogRepo struct {
	db *gorm.DB
}

func NewJobLogRepo(in any) *JobLogRepo {
	return &JobLogRepo{db: stores.GetCommonConn(in)}
}

type JobLogFilter struct {
	//todo 添加过滤字段
}

func (p JobLogRepo) fmtFilter(ctx context.Context, f JobLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p JobLogRepo) Insert(ctx context.Context, data *TimedJobLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p JobLogRepo) FindOneByFilter(ctx context.Context, f JobLogFilter) (*TimedJobLog, error) {
	var result TimedJobLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p JobLogRepo) FindByFilter(ctx context.Context, f JobLogFilter, page *def.PageInfo) ([]*TimedJobLog, error) {
	var results []*TimedJobLog
	db := p.fmtFilter(ctx, f).Model(&TimedJobLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p JobLogRepo) CountByFilter(ctx context.Context, f JobLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&TimedJobLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p JobLogRepo) Update(ctx context.Context, data *TimedJobLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p JobLogRepo) DeleteByFilter(ctx context.Context, f JobLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&TimedJobLog{}).Error
	return stores.ErrFmt(err)
}

func (p JobLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&TimedJobLog{}).Error
	return stores.ErrFmt(err)
}
func (p JobLogRepo) FindOne(ctx context.Context, id int64) (*TimedJobLog, error) {
	var result TimedJobLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p JobLogRepo) MultiInsert(ctx context.Context, data []*TimedJobLog) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&TimedJobLog{}).Create(data).Error
	return stores.ErrFmt(err)
}
