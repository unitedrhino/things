package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将ProtocolScript全局替换为模型的表名
2. 完善todo
*/

type ProtocolScriptRepo struct {
	db *gorm.DB
}

func NewProtocolScriptRepo(in any) *ProtocolScriptRepo {
	return &ProtocolScriptRepo{db: stores.GetCommonConn(in)}
}

type ProtocolScriptFilter struct {
	Name         string
	Status       int64
	TriggerDir   int64
	TriggerTimer int64
}

func (p ProtocolScriptRepo) fmtFilter(ctx context.Context, f ProtocolScriptFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.TriggerDir != 0 {
		db = db.Where("trigger_dir = ?", f.TriggerDir)
	}
	if f.TriggerTimer != 0 {
		db = db.Where("trigger_timer = ?", f.TriggerTimer)
	}
	return db
}

func (p ProtocolScriptRepo) Insert(ctx context.Context, data *DmProtocolScript) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProtocolScriptRepo) FindOneByFilter(ctx context.Context, f ProtocolScriptFilter) (*DmProtocolScript, error) {
	var result DmProtocolScript
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProtocolScriptRepo) FindByFilter(ctx context.Context, f ProtocolScriptFilter, page *stores.PageInfo) ([]*DmProtocolScript, error) {
	var results []*DmProtocolScript
	db := p.fmtFilter(ctx, f).Model(&DmProtocolScript{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProtocolScriptRepo) CountByFilter(ctx context.Context, f ProtocolScriptFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProtocolScript{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProtocolScriptRepo) Update(ctx context.Context, data *DmProtocolScript) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProtocolScriptRepo) DeleteByFilter(ctx context.Context, f ProtocolScriptFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProtocolScript{}).Error
	return stores.ErrFmt(err)
}

func (p ProtocolScriptRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProtocolScript{}).Error
	return stores.ErrFmt(err)
}
func (p ProtocolScriptRepo) FindOne(ctx context.Context, id int64) (*DmProtocolScript, error) {
	var result DmProtocolScript
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProtocolScriptRepo) MultiInsert(ctx context.Context, data []*DmProtocolScript) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProtocolScript{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d ProtocolScriptRepo) UpdateWithField(ctx context.Context, f ProtocolScriptFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProtocolScript{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
