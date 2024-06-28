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

type ProtocolInfoRepo struct {
	db *gorm.DB
}

func NewProtocolInfoRepo(in any) *ProtocolInfoRepo {
	return &ProtocolInfoRepo{db: stores.GetCommonConn(in)}
}

type ProtocolInfoFilter struct {
	ID            int64
	Name          string
	Code          string //  iThings,iThings-thingsboard,wumei,aliyun,huaweiyun,tuya
	TransProtocol string // 传输协议: mqtt,tcp,udp
	NotCodes      []string
}

func (p ProtocolInfoRepo) fmtFilter(ctx context.Context, f ProtocolInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.NotCodes) != 0 {
		db = db.Where("code not in ?", f.NotCodes)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	} else if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if f.TransProtocol != "" {
		db = db.Where("trans_protocol = ?", f.TransProtocol)
	}
	return db
}

func (p ProtocolInfoRepo) Insert(ctx context.Context, data *DmProtocolInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProtocolInfoRepo) FindOneByFilter(ctx context.Context, f ProtocolInfoFilter) (*DmProtocolInfo, error) {
	var result DmProtocolInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProtocolInfoRepo) FindByFilter(ctx context.Context, f ProtocolInfoFilter, page *stores.PageInfo) ([]*DmProtocolInfo, error) {
	var results []*DmProtocolInfo
	db := p.fmtFilter(ctx, f).Model(&DmProtocolInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProtocolInfoRepo) CountByFilter(ctx context.Context, f ProtocolInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProtocolInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProtocolInfoRepo) Update(ctx context.Context, data *DmProtocolInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProtocolInfoRepo) DeleteByFilter(ctx context.Context, f ProtocolInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProtocolInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ProtocolInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProtocolInfo{}).Error
	return stores.ErrFmt(err)
}
func (p ProtocolInfoRepo) FindOne(ctx context.Context, id int64) (*DmProtocolInfo, error) {
	var result DmProtocolInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProtocolInfoRepo) MultiInsert(ctx context.Context, data []*DmProtocolInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProtocolInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
