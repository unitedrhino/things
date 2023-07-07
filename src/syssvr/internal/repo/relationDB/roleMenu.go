package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleMenuRepo struct {
	db *gorm.DB
}

func NewRoleMenuRepo(in any) *RoleMenuRepo {
	return &RoleMenuRepo{db: store.GetCommonConn(in)}
}

type RoleMenuFilter struct {
	RoleIDs []int64
}

func (p RoleMenuRepo) fmtFilter(ctx context.Context, f RoleMenuFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.RoleIDs) != 0 {
		db = db.Where("`roleID` in ?", f.RoleIDs)
	}
	return db
}

func (g RoleMenuRepo) Insert(ctx context.Context, data *SysRoleMenu) error {
	result := g.db.WithContext(ctx).Create(data)
	return store.ErrFmt(result.Error)
}

func (g RoleMenuRepo) FindOneByFilter(ctx context.Context, f RoleMenuFilter) (*SysRoleMenu, error) {
	var result SysRoleMenu
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}
func (p RoleMenuRepo) FindByFilter(ctx context.Context, f RoleMenuFilter, page *def.PageInfo) ([]*SysRoleMenu, error) {
	var results []*SysRoleMenu
	db := p.fmtFilter(ctx, f).Model(&SysRoleMenu{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return results, nil
}

func (p RoleMenuRepo) CountByFilter(ctx context.Context, f RoleMenuFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysRoleMenu{})
	err = db.Count(&size).Error
	return size, store.ErrFmt(err)
}

func (g RoleMenuRepo) Update(ctx context.Context, data *SysRoleMenu) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return store.ErrFmt(err)
}

func (g RoleMenuRepo) DeleteByFilter(ctx context.Context, f RoleMenuFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysRoleMenu{}).Error
	return store.ErrFmt(err)
}

func (g RoleMenuRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&SysRoleMenu{}).Error
	return store.ErrFmt(err)
}
func (g RoleMenuRepo) FindOne(ctx context.Context, id int64) (*SysRoleMenu, error) {
	var result SysRoleMenu
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m RoleMenuRepo) MultiInsert(ctx context.Context, data []*SysRoleMenu) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysRoleMenu{}).Create(data).Error
	return store.ErrFmt(err)
}

func (g RoleMenuRepo) MultiUpdate(ctx context.Context, roleID int64, menuIDs []int64) error {
	var datas []*SysRoleMenu
	for _, v := range menuIDs {
		datas = append(datas, &SysRoleMenu{
			RoleID: roleID,
			MenuID: v,
		})
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		rm := NewRoleMenuRepo(tx)
		err := rm.DeleteByFilter(ctx, RoleMenuFilter{RoleIDs: []int64{roleID}})
		if err != nil {
			return err
		}
		err = rm.MultiInsert(ctx, datas)
		if err != nil {
			return err
		}
		return nil
	})
	return store.ErrFmt(err)
}
