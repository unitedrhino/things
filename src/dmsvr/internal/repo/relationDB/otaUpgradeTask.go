package relationDB

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type OtaUpgradeTaskRepo struct {
	db *gorm.DB
}

func NewOtaUpgradeTaskRepo(in any) *OtaUpgradeTaskRepo {
	return &OtaUpgradeTaskRepo{db: stores.GetCommonConn(in)}
}

type OtaUpgradeTaskFilter struct {
	Ids              []int64
	JobId            int64
	ProductId        string
	DeviceName       string
	DeviceNames      []string
	WithScheduleTime bool
	//TaskStatus     int64
	TaskStatusList []int
	ModuleName     string
}

func (p OtaUpgradeTaskRepo) fmtFilter(ctx context.Context, f OtaUpgradeTaskFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	if f.JobId != 0 {
		db = db.Where("job_id = ?", f.JobId)
	}
	if f.ProductId != "" {
		db = db.Where("product_id = ?", f.ProductId)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name like ?", "%"+f.DeviceName+"%")
	}
	if len(f.DeviceNames) != 0 {
		db = db.Where("device_name in ?", f.DeviceNames)
	}
	if f.WithScheduleTime {
		db = db.Where("schedule_time not null")
	}
	if len(f.TaskStatusList) != 0 {
		db = db.Where("task_status in ?", f.TaskStatusList)
	}
	if f.ModuleName != "" {
		db = db.Where("module_name = ?", f.ModuleName)
	}

	return db
}

func (p OtaUpgradeTaskRepo) Insert(ctx context.Context, data *DmOtaUpgradeTask) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OtaUpgradeTaskRepo) FindOneByFilter(ctx context.Context, f OtaUpgradeTaskFilter) (*DmOtaUpgradeTask, error) {
	var result DmOtaUpgradeTask
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaUpgradeTaskRepo) FindByFilter(ctx context.Context, f OtaUpgradeTaskFilter, page *def.PageInfo) ([]*DmOtaUpgradeTask, error) {
	var results []*DmOtaUpgradeTask
	db := p.fmtFilter(ctx, f).Model(&DmOtaUpgradeTask{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaUpgradeTaskRepo) CountByFilter(ctx context.Context, f OtaUpgradeTaskFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaUpgradeTask{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OtaUpgradeTaskRepo) Update(ctx context.Context, data *DmOtaUpgradeTask) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OtaUpgradeTaskRepo) DeleteByFilter(ctx context.Context, f OtaUpgradeTaskFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaUpgradeTask{}).Error
	return stores.ErrFmt(err)
}

func (p OtaUpgradeTaskRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DmOtaUpgradeTask{}).Error
	return stores.ErrFmt(err)
}
func (p OtaUpgradeTaskRepo) FindOne(ctx context.Context, id int64) (*DmOtaUpgradeTask, error) {
	var result DmOtaUpgradeTask
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OtaUpgradeTaskRepo) MultiInsert(ctx context.Context, data []*DmOtaUpgradeTask) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaUpgradeTask{}).Create(data).Error
	return stores.ErrFmt(err)
}

// 批量更新
func (p OtaUpgradeTaskRepo) BatchUpdateField(ctx context.Context, f OtaUpgradeTaskFilter, updateData map[string]interface{}) error {
	db := p.fmtFilter(ctx, f)
	err := db.WithContext(ctx).Model(&DmOtaUpgradeTask{}).Updates(updateData).Error
	return stores.ErrFmt(err)
}
