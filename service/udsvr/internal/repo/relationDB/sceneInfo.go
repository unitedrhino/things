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

type SceneInfoRepo struct {
	db *gorm.DB
}

func NewSceneInfoRepo(in any) *SceneInfoRepo {
	return &SceneInfoRepo{db: stores.GetCommonConn(in)}
}

type SceneInfoFilter struct {
	Name        string `json:"name"`
	Status      int64
	TriggerType string
	Tag         string
	AreaID      int64
}

func (p SceneInfoRepo) fmtFilter(ctx context.Context, f SceneInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//if f.AreaID != 0 {
	//	db = db.Where("area_id = ?", f.AreaID)
	//}
	if f.Tag != "" {
		db = db.Where("tag=?", f.Tag)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.TriggerType != "" {
		db = db.Where("trigger_type = ?", f.TriggerType)
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	return db
}

func (p SceneInfoRepo) Insert(ctx context.Context, data *UdSceneInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SceneInfoRepo) FindOneByFilter(ctx context.Context, f SceneInfoFilter) (*UdSceneInfo, error) {
	var result UdSceneInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SceneInfoRepo) FindByFilter(ctx context.Context, f SceneInfoFilter, page *def.PageInfo) ([]*UdSceneInfo, error) {
	var results []*UdSceneInfo
	db := p.fmtFilter(ctx, f).Model(&UdSceneInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SceneInfoRepo) CountByFilter(ctx context.Context, f SceneInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UdSceneInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SceneInfoRepo) Update(ctx context.Context, data *UdSceneInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SceneInfoRepo) DeleteByFilter(ctx context.Context, f SceneInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UdSceneInfo{}).Error
	return stores.ErrFmt(err)
}

func (p SceneInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UdSceneInfo{}).Error
	return stores.ErrFmt(err)
}
func (p SceneInfoRepo) FindOne(ctx context.Context, id int64) (*UdSceneInfo, error) {
	var result UdSceneInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SceneInfoRepo) MultiInsert(ctx context.Context, data []*UdSceneInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UdSceneInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
