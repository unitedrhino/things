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

type TaskLogRepo struct {
	db *gorm.DB
}

func NewJobLogRepo(in any) *TaskLogRepo {
	return &TaskLogRepo{db: stores.GetCommonConn(in)}
}

type TaskLogFilter struct {
	//todo 添加过滤字段
}

func (p TaskLogRepo) fmtFilter(ctx context.Context, f TaskLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p TaskLogRepo) Insert(ctx context.Context, data *TimedTaskLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TaskLogRepo) FindOneByFilter(ctx context.Context, f TaskLogFilter) (*TimedTaskLog, error) {
	var result TimedTaskLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TaskLogRepo) FindByFilter(ctx context.Context, f TaskLogFilter, page *def.PageInfo) ([]*TimedTaskLog, error) {
	var results []*TimedTaskLog
	db := p.fmtFilter(ctx, f).Model(&TimedTaskLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TaskLogRepo) CountByFilter(ctx context.Context, f TaskLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&TimedTaskLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TaskLogRepo) Update(ctx context.Context, data *TimedTaskLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TaskLogRepo) DeleteByFilter(ctx context.Context, f TaskLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&TimedTaskLog{}).Error
	return stores.ErrFmt(err)
}

func (p TaskLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&TimedTaskLog{}).Error
	return stores.ErrFmt(err)
}
func (p TaskLogRepo) FindOne(ctx context.Context, id int64) (*TimedTaskLog, error) {
	var result TimedTaskLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TaskLogRepo) MultiInsert(ctx context.Context, data []*TimedTaskLog) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&TimedTaskLog{}).Create(data).Error
	return stores.ErrFmt(err)
}
