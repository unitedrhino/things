package relationDB

import (
	"context"
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

type OtaJobRepo struct {
	db *gorm.DB
}

func NewOtaJobRepo(in any) *OtaJobRepo {
	return &OtaJobRepo{db: stores.GetCommonConn(in)}
}

type OtaJobFilter struct {
	//todo 添加过滤字段
	FirmwareID   int64
	ProductID    string
	DeviceName   string
	WithFirmware bool
	UpgradeType  int64
	Statues      []int64
	WithFiles    bool
}

func (p OtaJobRepo) fmtFilter(ctx context.Context, f OtaJobFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.WithFirmware {
		db = db.Preload("Firmware")
	}
	if f.ProductID != "" {
		db = db.Where("product_id=?", f.ProductID)
	}
	if f.FirmwareID != 0 {
		db = db.Where("firmware_id = ?", f.FirmwareID)
	}
	if len(f.Statues) != 0 {
		db = db.Where("status in ?", f.Statues)
	}
	if f.UpgradeType != 0 {
		db = db.Where("upgrade_type = ?", f.UpgradeType)
	}
	if f.WithFiles {
		db = db.Preload("Files")
	}
	return db
}

func (p OtaJobRepo) Insert(ctx context.Context, data *DmOtaFirmwareJob) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OtaJobRepo) FindOneByFilter(ctx context.Context, f OtaJobFilter) (*DmOtaFirmwareJob, error) {
	var result DmOtaFirmwareJob
	db := p.fmtFilter(ctx, f).Preload("Firmware")
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaJobRepo) FindByFilter(ctx context.Context, f OtaJobFilter, page *stores.PageInfo) ([]*DmOtaFirmwareJob, error) {
	var results []*DmOtaFirmwareJob
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareJob{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaJobRepo) CountByFilter(ctx context.Context, f OtaJobFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaFirmwareJob{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OtaJobRepo) Update(ctx context.Context, data *DmOtaFirmwareJob) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OtaJobRepo) DeleteByFilter(ctx context.Context, f OtaJobFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaFirmwareJob{}).Error
	return stores.ErrFmt(err)
}

func (p OtaJobRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmOtaFirmwareJob{}).Error
	return stores.ErrFmt(err)
}
func (p OtaJobRepo) FindOne(ctx context.Context, id int64) (*DmOtaFirmwareJob, error) {
	var result DmOtaFirmwareJob
	err := p.db.WithContext(ctx).Preload("Firmware").Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OtaJobRepo) MultiInsert(ctx context.Context, data []*DmOtaFirmwareJob) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaFirmwareJob{}).Create(data).Error
	return stores.ErrFmt(err)
}
