package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
)

type LoginLogRepo struct {
	db *gorm.DB
}

func NewLoginLogRepo(in any) *LoginLogRepo {
	return &LoginLogRepo{db: stores.GetCommonConn(in)}
}

type DateRange struct {
	Start string
	End   string
}
type LoginLogFilter struct {
	IpAddr        string
	LoginLocation string
	Data          *DateRange
}

func (p LoginLogRepo) fmtFilter(ctx context.Context, f LoginLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.IpAddr != "" {
		db = db.Where("`ipAddr`= ?", f.IpAddr)
	}
	if f.LoginLocation != "" {
		db = db.Where("`loginLocation` like ?", "%"+f.LoginLocation+"%")
	}
	if f.Data != nil && f.Data.Start != "" && f.Data.End != "" {
		db = db.Where("`createdTime` >= ? and `createdTime` <= ?", f.Data.Start, f.Data.End)
	}
	return db
}

func (g LoginLogRepo) Insert(ctx context.Context, data *SysLoginLog) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g LoginLogRepo) FindOneByFilter(ctx context.Context, f LoginLogFilter) (*SysLoginLog, error) {
	var result SysLoginLog
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p LoginLogRepo) FindByFilter(ctx context.Context, f LoginLogFilter, page *def.PageInfo) ([]*SysLoginLog, error) {
	var results []*SysLoginLog
	db := p.fmtFilter(ctx, f).Model(&SysLoginLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p LoginLogRepo) CountByFilter(ctx context.Context, f LoginLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysLoginLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g LoginLogRepo) Update(ctx context.Context, data *SysLoginLog) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g LoginLogRepo) DeleteByFilter(ctx context.Context, f LoginLogFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysLoginLog{}).Error
	return stores.ErrFmt(err)
}
func (g LoginLogRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&SysLoginLog{}).Error
	return stores.ErrFmt(err)
}

func (g LoginLogRepo) FindOne(ctx context.Context, id int64) (*SysLoginLog, error) {
	var result SysLoginLog
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
