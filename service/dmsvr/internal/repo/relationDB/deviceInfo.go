package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/stores"
	"gitee.com/i-Things/core/shared/utils"
	"gorm.io/gorm"
)

type DeviceInfoRepo struct {
	db *gorm.DB
}

type (
	DeviceFilter struct {
		ProductID         string
		AreaIDs           []int64
		DeviceName        string
		DeviceNames       []string
		Tags              map[string]string
		LastLoginTime     *def.TimeRange
		IsOnline          int64
		Range             int64
		Position          stores.Point
		DeviceAlias       string
		Versions          []string
		Cores             []*devices.Core
		WithProduct       bool
		ProductCategoryID int64
	}
)

func NewDeviceInfoRepo(in any) *DeviceInfoRepo {
	return &DeviceInfoRepo{db: stores.GetCommonConn(in)}
}

func (d DeviceInfoRepo) fmtFilter(ctx context.Context, f DeviceFilter) *gorm.DB {
	db := d.db.WithContext(ctx)
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	if f.ProductCategoryID != 0 {
		db = db.Where("product_id in (?)", db.Select("product_id").Model(DmProductInfo{}).Where("category_id=?", f.ProductCategoryID))
	}
	//业务过滤条件
	if f.ProductID != "" {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if len(f.AreaIDs) != 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	if len(f.Versions) != 0 {
		db = db.Where("version = ?", f.AreaIDs)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name like ?", "%"+f.DeviceName+"%")
	}
	if len(f.Cores) != 0 {
		scope := func(db *gorm.DB) *gorm.DB {
			for i, d := range f.Cores {
				if i == 0 {
					db = db.Where("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
					continue
				}
				db = db.Or("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
			}
			return db
		}
		db = db.Where(scope(db))
	}
	if len(f.DeviceNames) != 0 {
		db = db.Where("device_name in ?", f.DeviceNames)
	}
	if f.DeviceAlias != "" {
		db = db.Where("device_alias like ?", "%"+f.DeviceAlias+"%")
	}
	if f.Tags != nil {
		for k, v := range f.Tags {
			db = db.Where("JSON_CONTAINS(tags, JSON_OBJECT(?,?))",
				k, v)
		}
	}
	if f.LastLoginTime != nil {
		if f.LastLoginTime.Start != 0 {
			db = db.Where("last_login >= ?", utils.ToYYMMddHHSS(f.LastLoginTime.Start*1000))
		}
		if f.LastLoginTime.End != 0 {
			db = db.Where("last_login <= ?", utils.ToYYMMddHHSS(f.LastLoginTime.End*1000))
		}
	}

	if f.IsOnline != 0 {
		db = db.Where("is_online = ?", f.IsOnline)
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
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (d DeviceInfoRepo) Update(ctx context.Context, data *DmDeviceInfo) error {
	err := d.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d DeviceInfoRepo) Delete(ctx context.Context, id int64) error {
	err := d.db.WithContext(ctx).Where("id=?", id).Delete(&DmDeviceInfo{}).Error
	return stores.ErrFmt(err)
}

func (d DeviceInfoRepo) Insert(ctx context.Context, data *DmDeviceInfo) error {
	result := d.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (d DeviceInfoRepo) FindByFilter(ctx context.Context, f DeviceFilter, page *def.PageInfo) ([]*DmDeviceInfo, error) {
	var results []*DmDeviceInfo
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (d DeviceInfoRepo) CountByFilter(ctx context.Context, f DeviceFilter) (size int64, err error) {
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

type countModel struct {
	CountKey string
	Count    int64
}

func (d DeviceInfoRepo) CountGroupByField(ctx context.Context, f DeviceFilter, columnName string) (map[string]int64, error) {
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	countModelList := make([]*countModel, 0)
	err := db.Select(fmt.Sprintf("%s as CountKey", columnName), "count(1) as count").Group(columnName).Find(&countModelList).Error
	result := make(map[string]int64, 0)
	for _, v := range countModelList {
		result[v.CountKey] = v.Count
	}
	return result, stores.ErrFmt(err)
}

func (d DeviceInfoRepo) MultiUpdate(ctx context.Context, devices []*devices.Core, info *DmDeviceInfo) error {
	db := d.fmtFilter(ctx, DeviceFilter{Cores: devices}).Model(&DmDeviceInfo{})
	err := db.Updates(info).Error
	return stores.ErrFmt(err)
}
