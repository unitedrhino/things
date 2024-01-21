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
1. 将UserRole全局替换为模型的表名
2. 完善todo
*/

type UserRoleRepo struct {
	db *gorm.DB
}

func NewUserRoleRepo(in any) *UserRoleRepo {
	return &UserRoleRepo{db: stores.GetCommonConn(in)}
}

type UserRoleFilter struct {
	UserID int64
}

func (p UserRoleRepo) fmtFilter(ctx context.Context, f UserRoleFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UserID != 0 {
		db = db.Where("user_id =?", f.UserID)
	}
	return db
}

func (p UserRoleRepo) Insert(ctx context.Context, data *SysUserRole) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserRoleRepo) FindOneByFilter(ctx context.Context, f UserRoleFilter) (*SysUserRole, error) {
	var result SysUserRole
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserRoleRepo) FindByFilter(ctx context.Context, f UserRoleFilter, page *def.PageInfo) ([]*SysUserRole, error) {
	var results []*SysUserRole
	db := p.fmtFilter(ctx, f).Model(&SysUserRole{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserRoleRepo) CountByFilter(ctx context.Context, f UserRoleFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserRole{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserRoleRepo) Update(ctx context.Context, data *SysUserRole) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserRoleRepo) DeleteByFilter(ctx context.Context, f UserRoleFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserRole{}).Error
	return stores.ErrFmt(err)
}

func (p UserRoleRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysUserRole{}).Error
	return stores.ErrFmt(err)
}
func (p UserRoleRepo) FindOne(ctx context.Context, id int64) (*SysUserRole, error) {
	var result SysUserRole
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserRoleRepo) MultiInsert(ctx context.Context, data []*SysUserRole) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserRole{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p UserRoleRepo) MultiUpdate(ctx context.Context, userID int64, roleIDs []int64) error {
	var datas []*SysUserRole
	for _, v := range roleIDs {
		datas = append(datas, &SysUserRole{
			RoleID: v,
			UserID: userID,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewUserRoleRepo(tx)
		err := rm.DeleteByFilter(ctx, UserRoleFilter{UserID: userID})
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
