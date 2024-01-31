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

type OtaTaskDevicesRepo struct {
	db *gorm.DB
}

func NewOtaTaskDevicesRepo(in any) *OtaTaskDevicesRepo {
	return &OtaTaskDevicesRepo{db: stores.GetCommonConn(in)}
}

type OtaTaskDevicesFilter struct {
	FirmwareID int64
	TaskUid    string
	ProductID  string
	DeviceName string
	Status     int64
	Version    string
}

func (p OtaTaskDevicesRepo) fmtFilter(ctx context.Context, f OtaTaskDevicesFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.FirmwareID != 0 {
		db = db.Where("firmwareID=?", f.FirmwareID)
	}
	if f.TaskUid != "" {
		db = db.Where("taskUid=?", f.TaskUid)
	}
	if f.ProductID != "" {
		db = db.Where("productID=?", f.ProductID)
	}
	if f.DeviceName != "" {
		db = db.Where("deviceName=?", f.DeviceName)
	}
	if f.Version != "" {
		db = db.Where("version=?", f.Version)
	}
	if f.Status != 0 {
		db = db.Where("status=?", f.Status)
	}
	return db
}

func (g OtaTaskDevicesRepo) Insert(ctx context.Context, data *DmOtaTaskDevices) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g OtaTaskDevicesRepo) FindOneByFilter(ctx context.Context, f OtaTaskDevicesFilter) (*DmOtaTaskDevices, error) {
	var result DmOtaTaskDevices
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OtaTaskDevicesRepo) FindByFilter(ctx context.Context, f OtaTaskDevicesFilter, page *def.PageInfo) ([]*DmOtaTaskDevices, error) {
	var results []*DmOtaTaskDevices
	db := p.fmtFilter(ctx, f).Model(&DmOtaTaskDevices{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OtaTaskDevicesRepo) CountByFilter(ctx context.Context, f OtaTaskDevicesFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DmOtaTaskDevices{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g OtaTaskDevicesRepo) Update(ctx context.Context, data *DmOtaTaskDevices) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g OtaTaskDevicesRepo) DeleteByFilter(ctx context.Context, f OtaTaskDevicesFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&DmOtaTaskDevices{}).Error
	return stores.ErrFmt(err)
}

func (g OtaTaskDevicesRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("`id` = ?", id).Delete(&DmOtaTaskDevices{}).Error
	return stores.ErrFmt(err)
}
func (g OtaTaskDevicesRepo) FindOne(ctx context.Context, id int64) (*DmOtaTaskDevices, error) {
	var result DmOtaTaskDevices
	err := g.db.WithContext(ctx).Where("`id` = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (g OtaTaskDevicesRepo) CancelByTaskUid(ctx context.Context, taskUid string) error {
	err := g.db.Model(&DmOtaTaskDevices{}).WithContext(ctx).Where("taskUid = ? and status < 501", taskUid).Update("status", 701).Error
	return stores.ErrFmt(err)
}
func (o *OtaTaskDevicesFilter) GetEnabledBatchSql(sql *gorm.DB) *gorm.DB {
	sql = sql.Where("status < ?", 501)
	// version  怎么处理？先不处理了，返回给终端做判断，根据终端的返回设定当次升级的状态？
	return sql
}

// 查询设备当前可执行的升级批次
func (g OtaTaskDevicesRepo) FindEnableBatch(ctx context.Context, f OtaTaskDevicesFilter) (*DmOtaTaskDevices, error) {
	db := g.fmtFilter(ctx, f)
	f.GetEnabledBatchSql(db)
	var result *DmOtaTaskDevices
	qerr := db.Find(result).Error
	if qerr != nil {
		return nil, stores.ErrFmt(qerr)
	}
	return result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m OtaTaskDevicesRepo) MultiInsert(ctx context.Context, data []*DmOtaTaskDevices) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DmOtaTaskDevices{}).Create(data).Error
	return stores.ErrFmt(err)
}
