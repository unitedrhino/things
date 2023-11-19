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
1. 将Task全局替换为模型的表名
2. 完善todo
*/

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskInfoRepo(in any) *TaskRepo {
	return &TaskRepo{db: stores.GetCommonConn(in)}
}

type TaskFilter struct {
	IDs       []int64
	Types     []int64
	Status    []int64
	Codes     []string
	WithGroup bool
}

func (p TaskRepo) fmtFilter(ctx context.Context, f TaskFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.Status) != 0 {
		db = db.Where("status in ?", f.Status)
	}
	if len(f.Types) != 0 {
		db = db.Where("type in ?", f.Types)
	}
	if len(f.IDs) != 0 {
		db = db.Where("id in ?", f.IDs)
	}
	if len(f.Codes) > 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.WithGroup {
		db = db.Preload("Group")
	}
	return db
}

func (p TaskRepo) Insert(ctx context.Context, data *TimedTaskInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TaskRepo) FindOneByFilter(ctx context.Context, f TaskFilter) (*TimedTaskInfo, error) {
	var result TimedTaskInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TaskRepo) FindByFilter(ctx context.Context, f TaskFilter, page *def.PageInfo) ([]*TimedTaskInfo, error) {
	var results []*TimedTaskInfo
	db := p.fmtFilter(ctx, f).Model(&TimedTaskInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TaskRepo) CountByFilter(ctx context.Context, f TaskFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&TimedTaskInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TaskRepo) Update(ctx context.Context, data *TimedTaskInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TaskRepo) UpdateByFilter(ctx context.Context, data *TimedTaskInfo, f TaskFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Updates(data).Error
	return stores.ErrFmt(err)
}

func (p TaskRepo) DeleteByFilter(ctx context.Context, f TaskFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&TimedTaskInfo{}).Error
	return stores.ErrFmt(err)
}

func (p TaskRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&TimedTaskInfo{}).Error
	return stores.ErrFmt(err)
}
func (p TaskRepo) FindOne(ctx context.Context, id int64) (*TimedTaskInfo, error) {
	var result TimedTaskInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TaskRepo) MultiInsert(ctx context.Context, data []*TimedTaskInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&TimedTaskInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
