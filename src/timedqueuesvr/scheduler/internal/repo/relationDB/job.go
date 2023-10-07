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
1. 将Job全局替换为模型的表名
2. 完善todo
*/

type JobRepo struct {
	db *gorm.DB
}

func NewJobRepo(in any) *JobRepo {
	return &JobRepo{db: stores.GetCommonConn(in)}
}

type JobFilter struct {
	Status int64
}

func (p JobRepo) fmtFilter(ctx context.Context, f JobFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Status != 0 {
		db = db.Where("status=?", f.Status)
	}
	return db
}

func (p JobRepo) Insert(ctx context.Context, data *TimedQueueJob) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p JobRepo) FindOneByFilter(ctx context.Context, f JobFilter) (*TimedQueueJob, error) {
	var result TimedQueueJob
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p JobRepo) FindByFilter(ctx context.Context, f JobFilter, page *def.PageInfo) ([]*TimedQueueJob, error) {
	var results []*TimedQueueJob
	db := p.fmtFilter(ctx, f).Model(&TimedQueueJob{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p JobRepo) CountByFilter(ctx context.Context, f JobFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&TimedQueueJob{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p JobRepo) Update(ctx context.Context, data *TimedQueueJob) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p JobRepo) UpdateByFilter(ctx context.Context, data *TimedQueueJob, f JobFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Updates(data).Error
	return stores.ErrFmt(err)
}

func (p JobRepo) DeleteByFilter(ctx context.Context, f JobFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&TimedQueueJob{}).Error
	return stores.ErrFmt(err)
}

func (p JobRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&TimedQueueJob{}).Error
	return stores.ErrFmt(err)
}
func (p JobRepo) FindOne(ctx context.Context, id int64) (*TimedQueueJob, error) {
	var result TimedQueueJob
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p JobRepo) MultiInsert(ctx context.Context, data []*TimedQueueJob) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&TimedQueueJob{}).Create(data).Error
	return stores.ErrFmt(err)
}
