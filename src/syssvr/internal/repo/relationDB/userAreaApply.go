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

type UserAreaApplyRepo struct {
	db *gorm.DB
}

func NewUserAreaApplyRepo(in any) *UserAreaApplyRepo {
	return &UserAreaApplyRepo{db: stores.GetCommonConn(in)}
}

type UserAreaApplyFilter struct {
	ProjectID int64
	AuthTypes []int64
	IDs       []int64
	AreaIDs   []int64
}

func (p UserAreaApplyRepo) fmtFilter(ctx context.Context, f UserAreaApplyFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.AuthTypes) != 0 {
		db = db.Where("auth_type IN ?", f.AuthTypes)
	}
	if f.ProjectID != 0 {
		db = db.Where("project_id = ?", f.ProjectID)
	}
	if len(f.IDs) != 0 {
		db = db.Where("id IN ?", f.IDs)
	}
	if len(f.AreaIDs) != 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	return db
}

func (p UserAreaApplyRepo) Insert(ctx context.Context, data *SysUserAreaApply) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserAreaApplyRepo) FindOneByFilter(ctx context.Context, f UserAreaApplyFilter) (*SysUserAreaApply, error) {
	var result SysUserAreaApply
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserAreaApplyRepo) FindByFilter(ctx context.Context, f UserAreaApplyFilter, page *def.PageInfo) ([]*SysUserAreaApply, error) {
	var results []*SysUserAreaApply
	db := p.fmtFilter(ctx, f).Model(&SysUserAreaApply{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserAreaApplyRepo) CountByFilter(ctx context.Context, f UserAreaApplyFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserAreaApply{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserAreaApplyRepo) Update(ctx context.Context, data *SysUserAreaApply) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserAreaApplyRepo) DeleteByFilter(ctx context.Context, f UserAreaApplyFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserAreaApply{}).Error
	return stores.ErrFmt(err)
}

func (p UserAreaApplyRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysUserAreaApply{}).Error
	return stores.ErrFmt(err)
}
func (p UserAreaApplyRepo) FindOne(ctx context.Context, id int64) (*SysUserAreaApply, error) {
	var result SysUserAreaApply
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserAreaApplyRepo) MultiInsert(ctx context.Context, data []*SysUserAreaApply) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserAreaApply{}).Create(data).Error
	return stores.ErrFmt(err)
}
