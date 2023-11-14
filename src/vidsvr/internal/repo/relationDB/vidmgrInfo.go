package relationDB

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
)

type VidmgrInfoRepo struct {
	db *gorm.DB
}
type VidmgrConfigFilter struct {
	ApiSecret      string
	MediaServerIds []string
}

type VidmgrFilter struct {
	VidmgrType    int64
	VidmgrName    string
	VidmgrIDs     []string
	VidmgrNames   []string
	Tags          map[string]string
	LastLoginTime struct {
		Start int64
		End   int64
	}
	VidmgrStatus int64
}

type countModel struct {
	CountKey string
	Count    int64
}

func NewVidmgrtInfoRepo(in any) *VidmgrInfoRepo {
	return &VidmgrInfoRepo{db: stores.GetCommonConn(in)}
}

func (p VidmgrInfoRepo) fmtFilter(ctx context.Context, f VidmgrFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.VidmgrType != 0 {
		db = db.Where("type=?", f.VidmgrType)
	}
	if f.VidmgrName != "" {
		db = db.Where("name like ?", "%"+f.VidmgrName+"%")
	}
	if len(f.VidmgrIDs) != 0 {
		db = db.Where("id = ?", f.VidmgrIDs)
	}
	if len(f.VidmgrNames) != 0 {
		db = db.Where("name = ?", f.VidmgrNames)
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = db.Where("JSON_CONTAINS(tags, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	if f.VidmgrStatus != 0 {
		db = db.Where("status = ?", f.VidmgrStatus)
	}
	if f.LastLoginTime.Start != 0 {
		db = db.Where("last_login >= ?", utils.ToYYMMddHHSS(f.LastLoginTime.Start*1000))
	}
	if f.LastLoginTime.End != 0 {
		db = db.Where("last_login <= ?", utils.ToYYMMddHHSS(f.LastLoginTime.End*1000))
	}
	return db
}

func (p VidmgrInfoRepo) Insert(ctx context.Context, data *VidmgrInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p VidmgrInfoRepo) FindOneByFilter(ctx context.Context, f VidmgrFilter) (*VidmgrInfo, error) {
	var result VidmgrInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p VidmgrInfoRepo) Update(ctx context.Context, data *VidmgrInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.VidmgrID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p VidmgrInfoRepo) DeleteByFilter(ctx context.Context, f VidmgrFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&VidmgrInfo{}).Error
	return stores.ErrFmt(err)
}

func (p VidmgrInfoRepo) FindAllFilter(ctx context.Context, f VidmgrFilter) ([]*VidmgrInfo, error) {
	var results []*VidmgrInfo
	db := p.fmtFilter(ctx, f).Model(&VidmgrInfo{})
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrInfoRepo) FindByFilter(ctx context.Context, f VidmgrFilter, page *def.PageInfo) ([]*VidmgrInfo, error) {
	var results []*VidmgrInfo
	db := p.fmtFilter(ctx, f).Model(&VidmgrInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p VidmgrInfoRepo) CountByFilter(ctx context.Context, f VidmgrFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&VidmgrInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (d VidmgrInfoRepo) CountVidmgrByField(ctx context.Context, f VidmgrFilter, columnName string) (map[string]int64, error) {
	db := d.fmtFilter(ctx, f).Model(&VidmgrInfo{})
	countModelList := make([]*countModel, 0)
	err := db.Select(fmt.Sprintf("%s as CountKey", columnName), "count(1) as count").Group(columnName).Find(&countModelList).Error
	result := make(map[string]int64, 0)
	for _, v := range countModelList {
		result[v.CountKey] = v.Count
	}
	return result, stores.ErrFmt(err)
}
