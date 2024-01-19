package relationDB

import (
	"context"
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

type UserProjectRepo struct {
	db *gorm.DB
}

func NewUserProjectRepo(in any) *UserProjectRepo {
	return &UserProjectRepo{db: stores.GetCommonConn(in)}
}

type UserProjectFilter struct {
	UserID int64
}

func (p UserProjectRepo) fmtFilter(ctx context.Context, f UserProjectFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UserID != 0 {
		db = db.Where("user_id= ?", f.UserID)
	}
	return db
}

func (g UserProjectRepo) Insert(ctx context.Context, data *SysUserProject) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g UserProjectRepo) FindOneByFilter(ctx context.Context, f UserProjectFilter) (*SysUserProject, error) {
	var result SysUserProject
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserProjectRepo) FindByFilter(ctx context.Context, f UserProjectFilter, page *def.PageInfo) ([]*SysUserProject, error) {
	var results []*SysUserProject
	db := p.fmtFilter(ctx, f).Model(&SysUserProject{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserProjectRepo) CountByFilter(ctx context.Context, f UserProjectFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserProject{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g UserProjectRepo) Update(ctx context.Context, data *SysUserProject) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g UserProjectRepo) DeleteByFilter(ctx context.Context, f UserProjectFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysUserProject{}).Error
	return stores.ErrFmt(err)
}

func (g UserProjectRepo) Delete(ctx context.Context, userID int64, projectID int64) error {
	err := g.db.WithContext(ctx).Where("`userID` = ? and `projectID`=?", userID, projectID).
		Delete(&SysUserProject{}).Error
	return stores.ErrFmt(err)
}
func (g UserProjectRepo) FindOne(ctx context.Context, userID int64, projectID int64) (*SysUserProject, error) {
	var result SysUserProject
	err := g.db.WithContext(ctx).Where("`userID` = ? and `projectID`=?", userID, projectID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m UserProjectRepo) MultiInsert(ctx context.Context, data []*SysUserProject) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserProject{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (g UserProjectRepo) MultiUpdate(ctx context.Context, userID int64, projects []*userDataAuth.Project) error {
	var datas []*SysUserProject
	for _, v := range projects {
		datas = append(datas, &SysUserProject{
			UserID:    userID,
			ProjectID: v.ProjectID,
		})
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		rm := NewUserProjectRepo(tx)
		err := rm.DeleteByFilter(ctx, UserProjectFilter{UserID: userID})
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
