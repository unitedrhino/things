package relationDB

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
	"gorm.io/gorm"
)

type DeviceInfoRepo struct {
	db *gorm.DB
}

type (
	DeviceFilter struct {
		ProductID     string
		AreaIDs       []int64
		DeviceName    string
		DeviceNames   []string
		Tags          map[string]string
		LastLoginTime struct {
			Start int64
			End   int64
		}
		IsOnline    int64
		Range       int64
		Position    store.Point
		DeviceAlias string
	}
)

func NewDeviceInfoRepo(in any) *DeviceInfoRepo {
	return &DeviceInfoRepo{db: store.GetCommonConn(in)}
}

func (p DeviceInfoRepo) fmtFilter(ctx context.Context, f DeviceFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//业务过滤条件
	if f.ProductID != "" {
		db = db.Where("`productID` = ?", f.ProductID)
	}
	if len(f.AreaIDs) != 0 {
		db = db.Where(fmt.Sprintf("AreaID in (%v)", store.ArrayToSql(f.AreaIDs)))
	}
	if f.DeviceName != "" {
		db = db.Where("`deviceName` like ?", "%"+f.DeviceName+"%")
	}
	if len(f.DeviceNames) != 0 {
		db = db.Where(fmt.Sprintf("`deviceName` in (%v)", store.ArrayToSql(f.DeviceNames)))
	}
	if f.DeviceAlias != "" {
		db = db.Where("`deviceAlias` like ?", "%"+f.DeviceAlias+"%")
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = db.Where("JSON_CONTAINS(`tags`, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	if f.LastLoginTime.Start != 0 {
		db = db.Where("`lastLogin` >= ?", utils.ToYYMMddHHSS(f.LastLoginTime.Start*1000))
	}
	if f.LastLoginTime.End != 0 {
		db = db.Where("`lastLogin` <= ?", utils.ToYYMMddHHSS(f.LastLoginTime.End*1000))
	}
	if f.IsOnline != 0 {
		db = db.Where("`isOnline` = ?", f.IsOnline)
	}

	if f.Range > 0 {
		//f.Position 形如：point(116.393 39.905)
		db = db.Where(fmt.Sprintf(
			"round(st_distance_sphere(ST_GeomFromText(POINT(%v %v)), ST_GeomFromText(AsText(`position`))),2)>%d",
			f.Position.Longitude, f.Position.Latitude, f.Range))
	}
	return db
}

func (d DeviceInfoRepo) FindOneByFilter(ctx context.Context, f DeviceFilter) (*DmDeviceInfo, error) {
	var result DmDeviceInfo
	db := d.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return &result, nil
}

func (d DeviceInfoRepo) FindOneByIccid(ctx context.Context, iccid sql.NullString) (*DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindOneByPhone(ctx context.Context, phone sql.NullString) (*DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) FindOneByProductIDDeviceName(ctx context.Context, productID string, deviceName string) (*DmDeviceInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d DeviceInfoRepo) Update(ctx context.Context, data *DmDeviceInfo) error {
	err := d.db.WithContext(ctx).Where("`id` = ?", data.ID).Save(data).Error
	return store.ErrFmt(err)
}

func (d DeviceInfoRepo) Delete(ctx context.Context, id int64) error {
	err := d.db.WithContext(ctx).Where("`id`=?", id).Delete(&DmDeviceInfo{}).Error
	return store.ErrFmt(err)
}

func (d DeviceInfoRepo) Insert(ctx context.Context, data *DmDeviceInfo) error {
	result := d.db.WithContext(ctx).Create(data)
	return store.ErrFmt(result.Error)
}

func (d DeviceInfoRepo) FindByFilter(ctx context.Context, f DeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error) {
	var results []*DmDeviceInfo
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, store.ErrFmt(err)
	}
	return results, nil
}

func (d DeviceInfoRepo) CountByFilter(ctx context.Context, f DeviceFilter) (size int64, err error) {
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	err = db.Count(&size).Error
	return size, store.ErrFmt(err)
}

type countModel struct {
	CountKey string
	Count    int64
}

func (d DeviceInfoRepo) CountGroupByField(ctx context.Context, f DeviceFilter, columnName string) (map[string]int64, error) {
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	countModelList := make([]*countModel, 0)
	err := db.Select(fmt.Sprintf("`%s` as CountKey", columnName), "count(1) as count").Group(columnName).Find(&countModelList).Error
	result := make(map[string]int64, 0)
	for _, v := range countModelList {
		result[v.CountKey] = v.Count
	}
	return result, err
}

func (d DeviceInfoRepo) UpdateDeviceInfo(ctx context.Context, data *DmDeviceInfo) error {
	err := d.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return store.ErrFmt(err)
}
