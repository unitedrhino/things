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

type OtaFirmwareFileRepo struct {
	db *gorm.DB
}

func NewOtaFirmwareFileRepo(in any) *OtaFirmwareFileRepo {
	return &OtaFirmwareFileRepo{db: stores.GetCommonConn(in)}
}

type OtaFirmwareFileFilter struct {
	ID         int64
	FirmwareID int64
	Size       *int64
}

func (p OtaFirmwareFileRepo) fmtFilter(ctx context.Context, f OtaFirmwareFileFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.FirmwareID != 0 {
		db = db.Where("firmware_id=?", f.FirmwareID)
	}
	if f.Size != nil {
		db = db.Where("size=?", *f.Size)
	}
	return db
}

func (g OtaFirmwareFileRepo) Insert(ctx context.Context, data *DmOtaFirmwareFile) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g OtaFirmwareFileRepo) FindOneByFilter(ctx context.Context, f OtaFirmwareFileFilter) (*DmOtaFirmwareFile, error) {
	var result DmOtaFirmwareFile
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaFirmwareFileRepo) FindByFilter(ctx context.Context, f OtaFirmwareFileFilter, page *def.PageInfo) ([]*DmOtaFirmwareFile, error) {
	var results []*DmOtaFirmwareFile
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareFile{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaFirmwareFileRepo) CountByFilter(ctx context.Context, f OtaFirmwareFileFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareFile{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g OtaFirmwareFileRepo) Update(ctx context.Context, data *DmOtaFirmwareFile) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g OtaFirmwareFileRepo) DeleteByFilter(ctx context.Context, f OtaFirmwareFileFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaFirmwareFile{}).Error
	return stores.ErrFmt(err)
}

func (g OtaFirmwareFileRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&DmOtaFirmwareFile{}).Error
	return stores.ErrFmt(err)
}
func (g OtaFirmwareFileRepo) FindOne(ctx context.Context, id int64) (*DmOtaFirmwareFile, error) {
	var result DmOtaFirmwareFile
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m OtaFirmwareFileRepo) MultiInsert(ctx context.Context, data []*DmOtaFirmwareFile) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaFirmwareFile{}).Create(data).Error
	return stores.ErrFmt(err)
}
