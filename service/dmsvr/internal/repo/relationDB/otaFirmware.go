package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type OtaFirmwareInfoRepo struct {
	db *gorm.DB
}

func NewOtaFirmwareInfoRepo(in any) *OtaFirmwareInfoRepo {
	return &OtaFirmwareInfoRepo{db: stores.GetCommonConn(in)}
}

type OtaFirmwareInfoFilter struct {
	ProductID   string
	Name        string
	ID          int64
	Version     string
	WithProduct bool
	WithFiles   bool
}

func (p OtaFirmwareInfoRepo) fmtFilter(ctx context.Context, f OtaFirmwareInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.Version != "" {
		db = db.Where("version=?", f.Version)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.WithFiles {
		db = db.Preload("Files")
	}
	return db
}

func (g OtaFirmwareInfoRepo) Insert(ctx context.Context, data *DmOtaFirmwareInfo) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g OtaFirmwareInfoRepo) FindOneByFilter(ctx context.Context, f OtaFirmwareInfoFilter) (*DmOtaFirmwareInfo, error) {
	var result DmOtaFirmwareInfo
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaFirmwareInfoRepo) FindByFilter(ctx context.Context, f OtaFirmwareInfoFilter, page *def.PageInfo) ([]*DmOtaFirmwareInfo, error) {
	var results []*DmOtaFirmwareInfo
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaFirmwareInfoRepo) CountByFilter(ctx context.Context, f OtaFirmwareInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g OtaFirmwareInfoRepo) Update(ctx context.Context, data *DmOtaFirmwareInfo) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g OtaFirmwareInfoRepo) DeleteByFilter(ctx context.Context, f OtaFirmwareInfoFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaFirmwareInfo{}).Error
	return stores.ErrFmt(err)
}

func (g OtaFirmwareInfoRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&DmOtaFirmwareInfo{}).Error
	return stores.ErrFmt(err)
}
func (g OtaFirmwareInfoRepo) FindOne(ctx context.Context, id int64) (*DmOtaFirmwareInfo, error) {
	var result DmOtaFirmwareInfo
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m OtaFirmwareInfoRepo) MultiInsert(ctx context.Context, data []*DmOtaFirmwareInfo) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaFirmwareInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
