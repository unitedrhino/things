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
1. 将example全局替换为模型的表名
2. 完善todo
*/

type OtaFirmwareRepo struct {
	db *gorm.DB
}

func NewOtaFirmwareRepo(in any) *OtaFirmwareRepo {
	return &OtaFirmwareRepo{db: stores.GetCommonConn(in)}
}

type OtaFirmwareFilter struct {
	ProductID  string
	FirmwareID int64
	Version    string
}

func (p OtaFirmwareRepo) fmtFilter(ctx context.Context, f OtaFirmwareFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.FirmwareID != 0 {
		db = db.Where("id=?", f.FirmwareID)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.Version != "" {
		db = db.Where("version=?", f.Version)
	}
	return db
}

func (g OtaFirmwareRepo) Insert(ctx context.Context, data *DmOtaFirmware) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g OtaFirmwareRepo) FindOneByFilter(ctx context.Context, f OtaFirmwareFilter) (*DmOtaFirmware, error) {
	var result DmOtaFirmware
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaFirmwareRepo) FindByFilter(ctx context.Context, f OtaFirmwareFilter, page *def.PageInfo) ([]*DmOtaFirmware, error) {
	var results []*DmOtaFirmware
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmware{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaFirmwareRepo) CountByFilter(ctx context.Context, f OtaFirmwareFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmware{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g OtaFirmwareRepo) Update(ctx context.Context, data *DmOtaFirmware) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g OtaFirmwareRepo) DeleteByFilter(ctx context.Context, f OtaFirmwareFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaFirmware{}).Error
	return stores.ErrFmt(err)
}

func (g OtaFirmwareRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&DmOtaFirmware{}).Error
	return stores.ErrFmt(err)
}
func (g OtaFirmwareRepo) FindOne(ctx context.Context, id int64) (*DmOtaFirmware, error) {
	var result DmOtaFirmware
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m OtaFirmwareRepo) MultiInsert(ctx context.Context, data []*DmOtaFirmware) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaFirmware{}).Create(data).Error
	return stores.ErrFmt(err)
}
