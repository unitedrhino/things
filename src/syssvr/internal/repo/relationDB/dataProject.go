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

type DataProjectRepo struct {
	db *gorm.DB
}

func NewDataProjectRepo(in any) *DataProjectRepo {
	return &DataProjectRepo{db: stores.GetCommonConn(in)}
}

type UserProjectFilter struct {
	UserID int64
}

func (p DataProjectRepo) fmtFilter(ctx context.Context, f UserProjectFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//if f.UserID != 0 {
	//	db = db.Where("user_id= ?", f.UserID)
	//}
	return db
}

func (g DataProjectRepo) Insert(ctx context.Context, data *SysDataProject) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g DataProjectRepo) FindOneByFilter(ctx context.Context, f UserProjectFilter) (*SysDataProject, error) {
	var result SysDataProject
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DataProjectRepo) FindByFilter(ctx context.Context, f UserProjectFilter, page *def.PageInfo) ([]*SysDataProject, error) {
	var results []*SysDataProject
	db := p.fmtFilter(ctx, f).Model(&SysDataProject{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DataProjectRepo) CountByFilter(ctx context.Context, f UserProjectFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDataProject{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g DataProjectRepo) Update(ctx context.Context, data *SysDataProject) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g DataProjectRepo) DeleteByFilter(ctx context.Context, f UserProjectFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysDataProject{}).Error
	return stores.ErrFmt(err)
}

func (g DataProjectRepo) Delete(ctx context.Context, targetType string, targetID int64, projectID int64) error {
	err := g.db.WithContext(ctx).Where("target_type=? and target_id = ? and project_id=?", targetType, targetID, projectID).
		Delete(&SysDataProject{}).Error
	return stores.ErrFmt(err)
}
func (g DataProjectRepo) FindOne(ctx context.Context, targetType string, targetID int64, projectID int64) (*SysDataProject, error) {
	var result SysDataProject
	err := g.db.WithContext(ctx).Where("target_type=? and target_id = ? and project_id=?", targetType, targetID, projectID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m DataProjectRepo) MultiInsert(ctx context.Context, data []*SysDataProject) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDataProject{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (g DataProjectRepo) MultiUpdate(ctx context.Context, userID int64, projects []*userDataAuth.Project) error {
	var datas []*SysDataProject
	for _, v := range projects {
		datas = append(datas, &SysDataProject{
			TargetID:  userID,
			ProjectID: v.ProjectID,
		})
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		rm := NewDataProjectRepo(tx)
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
