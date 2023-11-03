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

type TaskGroupRepo struct {
	db *gorm.DB
}

func NewTaskGroupRepo(in any) *TaskGroupRepo {
	return &TaskGroupRepo{db: stores.GetCommonConn(in)}
}

type TaskGroupFilter struct {
	Code string
}

func (p TaskGroupRepo) fmtFilter(ctx context.Context, f TaskGroupFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	return db
}

func (p TaskGroupRepo) Insert(ctx context.Context, data *TimedTaskGroup) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TaskGroupRepo) FindOneByFilter(ctx context.Context, f TaskGroupFilter) (*TimedTaskGroup, error) {
	var result TimedTaskGroup
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TaskGroupRepo) FindByFilter(ctx context.Context, f TaskGroupFilter, page *def.PageInfo) ([]*TimedTaskGroup, error) {
	var results []*TimedTaskGroup
	db := p.fmtFilter(ctx, f).Model(&TimedTaskGroup{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TaskGroupRepo) CountByFilter(ctx context.Context, f TaskGroupFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&TimedTaskGroup{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TaskGroupRepo) Update(ctx context.Context, data *TimedTaskGroup) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TaskGroupRepo) DeleteByFilter(ctx context.Context, f TaskGroupFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&TimedTaskGroup{}).Error
	return stores.ErrFmt(err)
}

func (p TaskGroupRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&TimedTaskGroup{}).Error
	return stores.ErrFmt(err)
}
func (p TaskGroupRepo) FindOne(ctx context.Context, id int64) (*TimedTaskGroup, error) {
	var result TimedTaskGroup
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TaskGroupRepo) MultiInsert(ctx context.Context, data []*TimedTaskGroup) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&TimedTaskGroup{}).Create(data).Error
	return stores.ErrFmt(err)
}
