package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type ProtocolServiceRepo struct {
	db *gorm.DB
}

func NewProtocolServiceRepo(in any) *ProtocolServiceRepo {
	return &ProtocolServiceRepo{db: stores.GetCommonConn(in)}
}

type ProtocolServiceFilter struct {
	ID   int64
	Code string //  iThings,iThings-thingsboard,wumei,aliyun,huaweiyun,tuya
	IP   string
	Port int64
}

func (p ProtocolServiceRepo) fmtFilter(ctx context.Context, f ProtocolServiceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	} else if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if f.Port != 0 {
		db = db.Where("port = ?", f.Port)
	}
	if f.IP != "" {
		db = db.Where("ip = ?", f.IP)
	}
	return db
}

func (p ProtocolServiceRepo) Insert(ctx context.Context, data *DmProtocolService) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProtocolServiceRepo) FindOneByFilter(ctx context.Context, f ProtocolServiceFilter) (*DmProtocolService, error) {
	var result DmProtocolService
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProtocolServiceRepo) FindByFilter(ctx context.Context, f ProtocolServiceFilter, page *stores.PageInfo) ([]*DmProtocolService, error) {
	var results []*DmProtocolService
	db := p.fmtFilter(ctx, f).Model(&DmProtocolService{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProtocolServiceRepo) CountByFilter(ctx context.Context, f ProtocolServiceFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProtocolService{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProtocolServiceRepo) Update(ctx context.Context, data *DmProtocolService) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProtocolServiceRepo) DownStatus(ctx context.Context) error {
	err := p.db.WithContext(ctx).Model(DmProtocolService{}).Where("status = ? and updated_time<?", def.True, time.Now().Add(-time.Minute*5)).
		Update("status", def.False).Error
	return stores.ErrFmt(err)
}

func (p ProtocolServiceRepo) DeleteByFilter(ctx context.Context, f ProtocolServiceFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProtocolService{}).Error
	return stores.ErrFmt(err)
}

func (p ProtocolServiceRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProtocolService{}).Error
	return stores.ErrFmt(err)
}
func (p ProtocolServiceRepo) FindOne(ctx context.Context, id int64) (*DmProtocolService, error) {
	var result DmProtocolService
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProtocolServiceRepo) MultiInsert(ctx context.Context, data []*DmProtocolService) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProtocolService{}).Create(data).Error
	return stores.ErrFmt(err)
}
