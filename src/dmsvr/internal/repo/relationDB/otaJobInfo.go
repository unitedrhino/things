package relationDB

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/stores"
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
	FirmwareId int64
	ProductId  string
	DeviceName string
}

func (p OtaJobRepo) fmtFilter(ctx context.Context, f OtaJobFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.FirmwareId != 0 {
		db = db.Where("firmware_id = ?", f.FirmwareId)
	}
	return db
}

func (p OtaJobRepo) Insert(ctx context.Context, data *DmOtaJob) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OtaJobRepo) FindOneByFilter(ctx context.Context, f OtaJobFilter) (*DmOtaJob, error) {
	var result DmOtaJob
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaJobRepo) FindByFilter(ctx context.Context, f OtaJobFilter, page *def.PageInfo) ([]*DmOtaJob, error) {
	var results []*DmOtaJob
	db := p.fmtFilter(ctx, f).Model(&DmOtaJob{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaJobRepo) CountByFilter(ctx context.Context, f OtaJobFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaJob{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OtaJobRepo) Update(ctx context.Context, data *DmOtaJob) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OtaJobRepo) DeleteByFilter(ctx context.Context, f OtaJobFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaJob{}).Error
	return stores.ErrFmt(err)
}

func (p OtaJobRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmOtaJob{}).Error
	return stores.ErrFmt(err)
}
func (p OtaJobRepo) FindOne(ctx context.Context, id int64) (*DmOtaJob, error) {
	var result DmOtaJob
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OtaJobRepo) MultiInsert(ctx context.Context, data []*DmOtaJob) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaJob{}).Create(data).Error
	return stores.ErrFmt(err)
}
