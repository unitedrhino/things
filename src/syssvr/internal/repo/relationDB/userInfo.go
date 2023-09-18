package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type UserInfoRepo struct {
	db *gorm.DB
}

func NewUserInfoRepo(in any) *UserInfoRepo {
	return &UserInfoRepo{db: stores.GetCommonConn(in)}
}

type UserInfoFilter struct {
	UserNames []string
	UserName  string
	Phone     string
	Email     string
	Accounts  []string //账号查询 非模糊查询
	Wechat    string
}

func (p UserInfoRepo) accountsFilter(db *gorm.DB, accounts []string) *gorm.DB {
	db = db.Where("user_name in ?", accounts).
		Or("email in ?", accounts).
		Or("phone in ?", accounts)
	return db
}

func (p UserInfoRepo) fmtFilter(ctx context.Context, f UserInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.UserNames) != 0 {
		db = db.Where("user_name in ?", f.UserNames)
	}
	if len(f.Accounts) != 0 {
		db = p.accountsFilter(db, f.Accounts)
	}
	if f.UserName != "" {
		db = db.Where("user_name like ?", "%"+f.UserName+"%")
	}
	if f.Phone != "" {
		db = db.Where("phone like ?", "%"+f.Phone+"%")
	}
	if f.Email != "" {
		db = db.Where("email like ?", "%"+f.Email+"%")
	}
	if f.Wechat != "" {
		db = db.Where("wechat = ?", f.Wechat)
	}
	return db
}

func (p UserInfoRepo) Insert(ctx context.Context, data *SysUserInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserInfoRepo) FindOneByFilter(ctx context.Context, f UserInfoFilter) (*SysUserInfo, error) {
	var result SysUserInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserInfoRepo) FindByFilter(ctx context.Context, f UserInfoFilter, page *def.PageInfo) ([]*SysUserInfo, error) {
	var results []*SysUserInfo
	db := p.fmtFilter(ctx, f).Model(&SysUserInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserInfoRepo) CountByFilter(ctx context.Context, f UserInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserInfoRepo) Update(ctx context.Context, data *SysUserInfo) error {
	err := p.db.WithContext(ctx).Where("user_id = ?", data.UserID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserInfoRepo) DeleteByFilter(ctx context.Context, f UserInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserInfo{}).Error
	return stores.ErrFmt(err)
}

func (p UserInfoRepo) Delete(ctx context.Context, userID int64) error {
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&SysUserInfo{}).Error
	return stores.ErrFmt(err)
}
func (p UserInfoRepo) FindOne(ctx context.Context, userID int64) (*SysUserInfo, error) {
	var result SysUserInfo
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).First(&result).Error
	return &result, stores.ErrFmt(err)
}
