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

type ProjectInfoRepo struct {
	db *gorm.DB
}

func NewProjectInfoRepo(in any) *ProjectInfoRepo {
	return &ProjectInfoRepo{db: stores.GetCommonConn(in)}
}

type ProjectInfoFilter struct {
	//todo 添加过滤字段
}

func (p ProjectInfoRepo) fmtFilter(ctx context.Context, f ProjectInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p ProjectInfoRepo) Insert(ctx context.Context, data *ViewProjectInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProjectInfoRepo) FindOneByFilter(ctx context.Context, f ProjectInfoFilter) (*ViewProjectInfo, error) {
	var result ViewProjectInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProjectInfoRepo) FindByFilter(ctx context.Context, f ProjectInfoFilter, page *def.PageInfo) ([]*ViewProjectInfo, error) {
	var results []*ViewProjectInfo
	db := p.fmtFilter(ctx, f).Model(&ViewProjectInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProjectInfoRepo) CountByFilter(ctx context.Context, f ProjectInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&ViewProjectInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProjectInfoRepo) Update(ctx context.Context, data *ViewProjectInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProjectInfoRepo) DeleteByFilter(ctx context.Context, f ProjectInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&ViewProjectInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ProjectInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&ViewProjectInfo{}).Error
	return stores.ErrFmt(err)
}
func (p ProjectInfoRepo) FindOne(ctx context.Context, id int64) (*ViewProjectInfo, error) {
	var result ViewProjectInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProjectInfoRepo) MultiInsert(ctx context.Context, data []*ViewProjectInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&ViewProjectInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
