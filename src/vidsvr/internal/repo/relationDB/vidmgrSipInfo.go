package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

// 用于保存gb28181服务的ID信息，还有匹配对应的连接信息
type VidmgrSipInfoRepo struct {
	db *gorm.DB
}

func NewVidmgrSipInfoRepo(in any) *VidmgrSipInfoRepo {
	return &VidmgrSipInfoRepo{db: stores.GetCommonConn(in)}
}

type VidmgrSipInfoFilter struct {
	IDs      []int64
	Region   string
	Cid      string
	Did      string
	Lid      string
	VidmgrID string
}

func (p VidmgrSipInfoRepo) fmtFilter(ctx context.Context, f VidmgrSipInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if len(f.IDs) != 0 {
		db = db.Where("id in?", f.IDs)
	}
	if f.Region != "" {
		db = db.Where("region =?", f.Region)
	}
	if f.Cid != "" {
		db = db.Where("cid =?", f.Cid)
	}
	if f.Did != "" {
		db = db.Where("did =?", f.Did)
	}
	if f.Lid != "" {
		db = db.Where("lid =?", f.Lid)
	}
	if f.VidmgrID != "" {
		db = db.Where("vidmgr_id =?", f.VidmgrID)
	}
	return db
}

func (p VidmgrSipInfoRepo) Insert(ctx context.Context, data *VidmgrSipInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p VidmgrSipInfoRepo) FindOneByFilter(ctx context.Context, f VidmgrSipInfoFilter) (*VidmgrSipInfo, error) {
	var result VidmgrSipInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p VidmgrSipInfoRepo) Update(ctx context.Context, data *VidmgrSipInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p VidmgrSipInfoRepo) DeleteByFilter(ctx context.Context, f VidmgrSipInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&VidmgrSipInfo{}).Error
	return stores.ErrFmt(err)
}

func (p VidmgrSipInfoRepo) FindByFilter(ctx context.Context, f VidmgrSipInfoFilter, page *def.PageInfo) ([]*VidmgrSipInfo, error) {
	var results []*VidmgrSipInfo
	db := p.fmtFilter(ctx, f).Model(&VidmgrSipInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrSipInfoRepo) CountByFilter(ctx context.Context, f VidmgrSipInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&VidmgrSipInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}
