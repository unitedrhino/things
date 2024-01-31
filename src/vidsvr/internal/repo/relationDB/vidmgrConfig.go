package relationDB

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/stores"
	"gorm.io/gorm"
)

type VidmgrConfigRepo struct {
	db *gorm.DB
}

func NewVidmgrConfigRepo(in any) *VidmgrConfigRepo {
	return &VidmgrConfigRepo{db: stores.GetCommonConn(in)}
}

type VidmgrConfigFilter struct {
	ApiSecret string
	VidmgrIDs []string
	ConfigIDs []string
	//VidmgrIPv4 int64
	VidmgrPort int64
	Secret     string
}

func (p VidmgrConfigRepo) fmtFilter(ctx context.Context, f VidmgrConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//if f.VidmgrIPv4 != 0 {
	//	db = db.Where("ipv4=?", f.VidmgrIPv4)
	//}
	//if f.VidmgrPort != 0 {
	//	db = db.Where("port=?", f.VidmgrPort)
	//}
	if f.Secret != "" {
		db = db.Where("secret=?", f.Secret)
	}
	if f.ApiSecret != "" {
		db = db.Where("secret=?", f.ApiSecret)
	}
	if len(f.VidmgrIDs) != 0 {
		db = db.Where("vidmgr_id in?", f.VidmgrIDs)
	}

	//if len(f.ConfigIDs) != 0 {
	//	db = db.Where("config_id=?", f.ConfigIDs)
	//}
	return db
}

func (p VidmgrConfigRepo) Insert(ctx context.Context, data *VidmgrConfig) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p VidmgrConfigRepo) FindOneByFilter(ctx context.Context, f VidmgrConfigFilter) (*VidmgrConfig, error) {
	var result VidmgrConfig
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p VidmgrConfigRepo) Update(ctx context.Context, data *VidmgrConfig) error {
	err := p.db.WithContext(ctx).Where("vidmgr_id = ?", data.VidmgrID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p VidmgrConfigRepo) DeleteByFilter(ctx context.Context, f VidmgrConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&VidmgrConfig{}).Error
	return stores.ErrFmt(err)
}

func (p VidmgrConfigRepo) FindByFilter(ctx context.Context, f VidmgrConfigFilter, page *def.PageInfo) ([]*VidmgrConfig, error) {
	var results []*VidmgrConfig
	db := p.fmtFilter(ctx, f).Model(&VidmgrConfig{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrConfigRepo) CountByFilter(ctx context.Context, f VidmgrConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&VidmgrConfig{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
