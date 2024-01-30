package relationDB

import (
	"context"
	"fmt"
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

type AreaInfoRepo struct {
	db *gorm.DB
}

func NewAreaInfoRepo(in any) *AreaInfoRepo {
	return &AreaInfoRepo{db: stores.GetCommonConn(in)}
}

type AreaInfoFilter struct {
	ProjectID    int64
	ParentAreaID int64
	AreaIDs      []int64
	AreaIDPath   string
	*AreaInfoWith
}

type AreaInfoWith struct {
	Children bool
	Parent   bool
}

func (p AreaInfoRepo) With(db *gorm.DB, with *AreaInfoWith) *gorm.DB {
	if with == nil {
		return db
	}
	if with.Parent {
		db = db.Preload("Parent")
	}
	if with.Children {
		db = db.Preload("Children")
	}
	return db
}

func (p AreaInfoRepo) fmtFilter(ctx context.Context, f AreaInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = p.With(db, f.AreaInfoWith)
	if f.ProjectID != 0 {
		db = db.Where("project_id = ?", f.ProjectID)
		//ctxs.SetMetaProjectID(ctx, f.ProjectID) //指定项目id的时候需要清除项目id
	}
	if f.ParentAreaID != 0 {
		db = db.Where("parent_area_id = ?", f.ParentAreaID)
	}
	if len(f.AreaIDs) != 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	if f.AreaIDPath != "" {
		db = db.Where("area_id_path like ?", f.AreaIDPath+"%")
	}
	return db
}

func (g AreaInfoRepo) Insert(ctx context.Context, data *SysAreaInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g AreaInfoRepo) FindOneByFilter(ctx context.Context, f AreaInfoFilter) (*SysAreaInfo, error) {
	var result SysAreaInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AreaInfoRepo) FindByFilter(ctx context.Context, f AreaInfoFilter, page *def.PageInfo) ([]*SysAreaInfo, error) {
	var results []*SysAreaInfo
	db := p.fmtFilter(ctx, f).Model(&SysAreaInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AreaInfoRepo) CountByFilter(ctx context.Context, f AreaInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysAreaInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g AreaInfoRepo) Update(ctx context.Context, data *SysAreaInfo) error {
	err := g.db.WithContext(ctx).Where("area_id = ?", data.AreaID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g AreaInfoRepo) DeleteByFilter(ctx context.Context, f AreaInfoFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysAreaInfo{}).Error
	return stores.ErrFmt(err)
}

func (g AreaInfoRepo) Delete(ctx context.Context, areaID int64) error {
	err := g.db.WithContext(ctx).Where("area_id = ?", areaID).Delete(&SysAreaInfo{}).Error
	return stores.ErrFmt(err)
}
func (g AreaInfoRepo) FindOne(ctx context.Context, areaID int64, with *AreaInfoWith) (*SysAreaInfo, error) {
	var result SysAreaInfo
	db := g.db.WithContext(ctx)
	db = g.With(db, with)
	err := db.Where("area_id = ?", areaID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (g AreaInfoRepo) FindIDsWithChildren(ctx context.Context, areaIDs []int64) ([]int64, error) {
	var resp []int64
	table := SysAreaInfo{}
	db := g.db.WithContext(ctx)
	sql := fmt.Sprintf(`
		WITH RECURSIVE cte AS (
		  SELECT area_id
		  FROM %s
		  WHERE area_id in ?
		  UNION ALL
		  SELECT t.area_id
		  FROM %s t
		  INNER JOIN cte ON t.parent_area_id = cte.area_id
		)
		SELECT area_id
		FROM cte`,
		table.TableName(), table.TableName())

	err := db.Raw(sql, areaIDs).Scan(&resp).Error
	return resp, stores.ErrFmt(err)

}

// 批量插入 LightStrategyDevice 记录
func (m AreaInfoRepo) MultiInsert(ctx context.Context, data []*SysAreaInfo) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysAreaInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
