package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectInfoRepo struct {
	db *gorm.DB
}

func NewProjectInfoRepo(in any) *ProjectInfoRepo {
	return &ProjectInfoRepo{db: stores.GetCommonConn(in)}
}

type ProjectInfoFilter struct {
	ProjectIDs  []int64
	ProjectName string
}

func (p ProjectInfoRepo) fmtFilter(ctx context.Context, f ProjectInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ProjectName != "" {
		db = db.Where("project_name like ?", "%"+f.ProjectName+"%")
	}
	if len(f.ProjectIDs) != 0 {
		db = db.Where("project_id in ?", f.ProjectIDs)
	}
	return db
}

func (g ProjectInfoRepo) Insert(ctx context.Context, data *SysProjectInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g ProjectInfoRepo) FindOneByFilter(ctx context.Context, f ProjectInfoFilter) (*SysProjectInfo, error) {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	var result SysProjectInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProjectInfoRepo) FindByFilter(ctx context.Context, f ProjectInfoFilter, page *def.PageInfo) ([]*SysProjectInfo, error) {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	var results []*SysProjectInfo
	db := p.fmtFilter(ctx, f).Model(&SysProjectInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProjectInfoRepo) CountByFilter(ctx context.Context, f ProjectInfoFilter) (size int64, err error) {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	db := p.fmtFilter(ctx, f).Model(&SysProjectInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g ProjectInfoRepo) Update(ctx context.Context, data *SysProjectInfo) error {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	err := g.db.WithContext(ctx).Where("project_id = ?", data.ProjectID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g ProjectInfoRepo) DeleteByFilter(ctx context.Context, f ProjectInfoFilter) error {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysProjectInfo{}).Error
	return stores.ErrFmt(err)
}

func (g ProjectInfoRepo) Delete(ctx context.Context, projectID int64) error {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	err := g.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&SysProjectInfo{}).Error
	return stores.ErrFmt(err)
}

func (g ProjectInfoRepo) FindOne(ctx context.Context, projectID int64) (*SysProjectInfo, error) {
	ctxs.ClearMetaProjectID(ctx) //默认情况下只返回当前项目,需要清除当前项目
	var result SysProjectInfo
	err := g.db.WithContext(ctx).Where("project_id = ?", projectID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m ProjectInfoRepo) MultiInsert(ctx context.Context, data []*SysProjectInfo) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysProjectInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
