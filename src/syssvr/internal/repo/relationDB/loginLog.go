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
		db = db.Where("ip_addr= ?", f.IpAddr)
	}
	if f.LoginLocation != "" {
		db = db.Where("login_location like ?", "%"+f.LoginLocation+"%")
	}
	if f.Data != nil && f.Data.Start != "" && f.Data.End != "" {
		db = db.Where("created_time >= ? and created_time <= ?", f.Data.Start, f.Data.End)
	}
	return db
}

func (p LoginLogRepo) Insert(ctx context.Context, data *SysTenantLoginLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p LoginLogRepo) FindOneByFilter(ctx context.Context, f LoginLogFilter) (*SysTenantLoginLog, error) {
	var result SysTenantLoginLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p LoginLogRepo) FindByFilter(ctx context.Context, f LoginLogFilter, page *def.PageInfo) ([]*SysTenantLoginLog, error) {
	var results []*SysTenantLoginLog
	db := p.fmtFilter(ctx, f).Model(&SysTenantLoginLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p LoginLogRepo) CountByFilter(ctx context.Context, f LoginLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantLoginLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p LoginLogRepo) Update(ctx context.Context, data *SysTenantLoginLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p LoginLogRepo) DeleteByFilter(ctx context.Context, f LoginLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantLoginLog{}).Error
	return stores.ErrFmt(err)
}
func (p LoginLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantLoginLog{}).Error
	return stores.ErrFmt(err)
}

func (p LoginLogRepo) FindOne(ctx context.Context, id int64) (*SysTenantLoginLog, error) {
	var result SysTenantLoginLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
