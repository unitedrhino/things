package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/userDataAuth"
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

type UserAreaRepo struct {
	db *gorm.DB
}

func NewUserAreaRepo(in any) *UserAreaRepo {
	return &UserAreaRepo{db: stores.GetCommonConn(in)}
}

type UserAreaFilter struct {
	UserID    int64
	ProjectID int64
	AreaIDs   []int64
}

func (p UserAreaRepo) fmtFilter(ctx context.Context, f UserAreaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.AreaIDs) > 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	if f.UserID != 0 {
		db = db.Where("user_id= ?", f.UserID)
	}
	if f.ProjectID != 0 {
		db = db.Where("project_id= ?", f.ProjectID)
		ctxs.SetMetaProjectID(ctx, f.ProjectID) //指定项目id的时候需要清除项目id
	}
	return db
}

func (g UserAreaRepo) Insert(ctx context.Context, data *SysUserArea) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g UserAreaRepo) FindOneByFilter(ctx context.Context, f UserAreaFilter) (*SysUserArea, error) {
	var result SysUserArea
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserAreaRepo) FindByFilter(ctx context.Context, f UserAreaFilter, page *def.PageInfo) ([]*SysUserArea, error) {
	var results []*SysUserArea
	db := p.fmtFilter(ctx, f).Model(&SysUserArea{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserAreaRepo) CountByFilter(ctx context.Context, f UserAreaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserArea{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g UserAreaRepo) Update(ctx context.Context, data *SysUserArea) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g UserAreaRepo) DeleteByFilter(ctx context.Context, f UserAreaFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysUserArea{}).Error
	return stores.ErrFmt(err)
}

func (g UserAreaRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&SysUserArea{}).Error
	return stores.ErrFmt(err)
}
func (g UserAreaRepo) FindOne(ctx context.Context, id int64) (*SysUserArea, error) {
	var result SysUserArea
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m UserAreaRepo) MultiInsert(ctx context.Context, data []*SysUserArea) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserArea{}).Create(data).Error
	return stores.ErrFmt(err)
}
func (g UserAreaRepo) MultiUpdate(ctx context.Context, userID, projectID int64, areas []*userDataAuth.Area) error {
	var datas []*SysUserArea
	for _, v := range areas {
		datas = append(datas, &SysUserArea{
			UserID:    userID,
			ProjectID: projectID,
			AreaID:    v.AreaID,
		})
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		rm := NewUserAreaRepo(tx)
		err := rm.DeleteByFilter(ctx, UserAreaFilter{UserID: userID, ProjectID: projectID})
		if err != nil {
			return err
		}
		if len(datas) != 0 {
			err = rm.MultiInsert(ctx, datas)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}
