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

type ProjectDetailRepo struct {
	db *gorm.DB
}

func NewProjectDetailRepo(in any) *ProjectDetailRepo {
	return &ProjectDetailRepo{db: stores.GetCommonConn(in)}
}

type ProjectDetailFilter struct {
	//todo 添加过滤字段
}

func (p ProjectDetailRepo) fmtFilter(ctx context.Context, f ProjectDetailFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p ProjectDetailRepo) Insert(ctx context.Context, data *ViewProjectDetail) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProjectDetailRepo) FindOneByFilter(ctx context.Context, f ProjectDetailFilter) (*ViewProjectDetail, error) {
	var result ViewProjectDetail
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProjectDetailRepo) FindByFilter(ctx context.Context, f ProjectDetailFilter, page *def.PageInfo) ([]*ViewProjectDetail, error) {
	var results []*ViewProjectDetail
	db := p.fmtFilter(ctx, f).Model(&ViewProjectDetail{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProjectDetailRepo) CountByFilter(ctx context.Context, f ProjectDetailFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&ViewProjectDetail{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProjectDetailRepo) Update(ctx context.Context, data *ViewProjectDetail) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProjectDetailRepo) DeleteByFilter(ctx context.Context, f ProjectDetailFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&ViewProjectDetail{}).Error
	return stores.ErrFmt(err)
}

func (p ProjectDetailRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("project_id = ?", id).Delete(&ViewProjectDetail{}).Error
	return stores.ErrFmt(err)
}
func (p ProjectDetailRepo) FindOne(ctx context.Context, id int64) (*ViewProjectDetail, error) {
	var result ViewProjectDetail
	err := p.db.WithContext(ctx).Where("project_id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProjectDetailRepo) MultiInsert(ctx context.Context, data []*ViewProjectDetail) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&ViewProjectDetail{}).Create(data).Error
	return stores.ErrFmt(err)
}
