package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type VidmgrConfigRepo struct {
	db *gorm.DB
}

func NewVidmgrtConfigRepo(in any) *VidmgrConfigRepo {
	return &VidmgrConfigRepo{db: stores.GetCommonConn(in)}
}

func (p VidmgrConfigRepo) fmtFilter(ctx context.Context, f VidmgrConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ApiSecret != "" {
		db = db.Where("api_secret=?", f.ApiSecret)
	}
	if len(f.MediaServerIds) != 0 {
		db = db.Where("general_mediaServerId=?", f.MediaServerIds)
	}
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
	err := p.db.WithContext(ctx).Where("GeneralMediaServerId = ?", data.GeneralMediaServerId).Save(data).Error
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
