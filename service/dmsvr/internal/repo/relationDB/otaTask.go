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

type OtaTaskRepo struct {
	db *gorm.DB
}

func NewOtaTaskRepo(in any) *OtaTaskRepo {
	return &OtaTaskRepo{db: stores.GetCommonConn(in)}
}

type OtaTaskFilter struct {
	FirmwareID int64
	TaskUid    string
	Status     int64
}

func (p OtaTaskRepo) fmtFilter(ctx context.Context, f OtaTaskFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.FirmwareID != 0 {
		db = db.Where("firmwareID=?", f.FirmwareID)
	}
	if f.TaskUid != "" {
		db = db.Where("taskUid=?", f.TaskUid)
	}
	if f.Status != 0 {
		db = db.Where("status=?", f.Status)
	}
	return db
}

func (g OtaTaskRepo) Insert(ctx context.Context, data *DmOtaTask) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g OtaTaskRepo) FindOneByFilter(ctx context.Context, f OtaTaskFilter) (*DmOtaTask, error) {
	var result DmOtaTask
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaTaskRepo) FindByFilter(ctx context.Context, f OtaTaskFilter, page *def.PageInfo) ([]*DmOtaTask, error) {
	var results []*DmOtaTask
	db := p.fmtFilter(ctx, f).Model(&DmOtaTask{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaTaskRepo) CountByFilter(ctx context.Context, f OtaTaskFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaTask{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g OtaTaskRepo) Update(ctx context.Context, data *DmOtaTask) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g OtaTaskRepo) DeleteByFilter(ctx context.Context, f OtaTaskFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaTask{}).Error
	return stores.ErrFmt(err)
}

func (g OtaTaskRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&DmOtaTask{}).Error
	return stores.ErrFmt(err)
}
func (g OtaTaskRepo) FindOne(ctx context.Context, id int64) (*DmOtaTask, error) {
	var result DmOtaTask
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m OtaTaskRepo) MultiInsert(ctx context.Context, data []*DmOtaTask) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaTask{}).Create(data).Error
	return stores.ErrFmt(err)
}
