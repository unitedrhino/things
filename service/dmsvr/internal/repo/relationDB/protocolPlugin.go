package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/share/devices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将ProtocolPlugin全局替换为模型的表名
2. 完善todo
*/

type ProtocolPluginRepo struct {
	db *gorm.DB
}

func NewProtocolPluginRepo(in any) *ProtocolPluginRepo {
	return &ProtocolPluginRepo{db: stores.GetCommonConn(in)}
}

type ProtocolPluginFilter struct {
	Name          string
	TriggerSrc    int64
	TriggerDir    int64
	TriggerTimer  int64
	BindProductID string
	BindDevice    *devices.Core
}

func (p ProtocolPluginRepo) fmtFilter(ctx context.Context, f ProtocolPluginFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("Name like ?", "%"+f.Name+"%")
	}
	if f.TriggerSrc != 0 {
		db = db.Where("trigger_src = ?", f.TriggerSrc)
	}
	if f.TriggerDir != 0 {
		db = db.Where("trigger_dir = ?", f.TriggerDir)
	}
	if f.TriggerTimer != 0 {
		db = db.Where("trigger_timer = ?", f.TriggerTimer)
	}
	//if f.BindProductID != 0 {
	//	db = db.Where("BindProductID = ?", f.BindProductID)
	//}
	//if f.BindDevice != 0 {
	//	db = db.Where("BindDevice = ?", f.BindDevice)
	//}
	return db
}

func (p ProtocolPluginRepo) Insert(ctx context.Context, data *DmProtocolPlugin) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProtocolPluginRepo) FindOneByFilter(ctx context.Context, f ProtocolPluginFilter) (*DmProtocolPlugin, error) {
	var result DmProtocolPlugin
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProtocolPluginRepo) FindByFilter(ctx context.Context, f ProtocolPluginFilter, page *stores.PageInfo) ([]*DmProtocolPlugin, error) {
	var results []*DmProtocolPlugin
	db := p.fmtFilter(ctx, f).Model(&DmProtocolPlugin{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProtocolPluginRepo) CountByFilter(ctx context.Context, f ProtocolPluginFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmProtocolPlugin{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProtocolPluginRepo) Update(ctx context.Context, data *DmProtocolPlugin) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProtocolPluginRepo) DeleteByFilter(ctx context.Context, f ProtocolPluginFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmProtocolPlugin{}).Error
	return stores.ErrFmt(err)
}

func (p ProtocolPluginRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmProtocolPlugin{}).Error
	return stores.ErrFmt(err)
}
func (p ProtocolPluginRepo) FindOne(ctx context.Context, id int64) (*DmProtocolPlugin, error) {
	var result DmProtocolPlugin
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProtocolPluginRepo) MultiInsert(ctx context.Context, data []*DmProtocolPlugin) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmProtocolPlugin{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d ProtocolPluginRepo) UpdateWithField(ctx context.Context, f ProtocolPluginFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmProtocolPlugin{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
