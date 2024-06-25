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

type SceneLogRepo struct {
	db *gorm.DB
}

func NewSceneLogRepo(in any) *SceneLogRepo {
	return &SceneLogRepo{db: stores.GetCommonConn(in)}
}

type SceneLogFilter struct {
	Time    *def.TimeRange
	Status  int64
	SceneID int64
}

func (p SceneLogRepo) fmtFilter(ctx context.Context, f SceneLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.Time.ToGorm(db, "created_time")
	if f.SceneID != 0 {
		db = db.Where("scene_id=?", f.SceneID)
	}
	if f.Status != 0 {
		db = db.Where("status=?", f.Status)
	}
	return db
}

func (p SceneLogRepo) Insert(ctx context.Context, data *UdSceneLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SceneLogRepo) FindOneByFilter(ctx context.Context, f SceneLogFilter) (*UdSceneLog, error) {
	var result UdSceneLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneLogRepo) FindByFilter(ctx context.Context, f SceneLogFilter, page *def.PageInfo) ([]*UdSceneLog, error) {
	var results []*UdSceneLog
	db := p.fmtFilter(ctx, f).Model(&UdSceneLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SceneLogRepo) CountByFilter(ctx context.Context, f SceneLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdSceneLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneLogRepo) Update(ctx context.Context, data *UdSceneLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SceneLogRepo) DeleteByFilter(ctx context.Context, f SceneLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdSceneLog{}).Error
	return stores.ErrFmt(err)
}

func (p SceneLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdSceneLog{}).Error
	return stores.ErrFmt(err)
}
func (p SceneLogRepo) FindOne(ctx context.Context, id int64) (*UdSceneLog, error) {
	var result UdSceneLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneLogRepo) MultiInsert(ctx context.Context, data []*UdSceneLog) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdSceneLog{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d SceneLogRepo) UpdateWithField(ctx context.Context, f SceneLogFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UdSceneLog{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
