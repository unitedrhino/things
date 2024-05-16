package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
)

type DeviceInfoRepo struct {
	db *gorm.DB
}

type (
	DeviceFilter struct {
		TenantCode        string
		TenantCodes       []string
		ProjectIDs        []int64
		ProductID         string
		AreaIDs           []int64
		NotAreaIDs        []int64
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
		Gateway           *devices.Core
		WithProduct       bool
		ProductCategoryID int64
		SharedDevices     []*devices.Core
		SharedType        def.SelectType
		WithManufacturer  bool
		DeviceType        int64
		GroupID           int64
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
	if f.WithManufacturer {
		db = db.Preload("Manufacturer")
	}
	if len(f.TenantCodes) != 0 {
		db = db.Where("tenant_code in ?", f.TenantCodes)
	}
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	if f.ProductCategoryID != 0 {
		db = db.Where("product_id in (?)", db.Select("product_id").Model(DmProductInfo{}).Where("category_id=?", f.ProductCategoryID))
	}
	//业务过滤条件
	if f.ProductID != "" {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if len(f.NotAreaIDs) != 0 {
		db = db.Where("area_id not in ?", f.NotAreaIDs)
	}
	if len(f.AreaIDs) != 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	if len(f.Versions) != 0 {
		db = db.Where("version in ?", f.Versions)
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
		db = db.Where(f.Position.Range("position", f.Range))
	}
	if f.SharedType != 0 && len(f.SharedDevices) > 0 && ctxs.GetUserCtx(ctx) != nil { //如果要获取共享设备
		scope := func(db *gorm.DB) *gorm.DB {
			for i, d := range f.SharedDevices {
				if i == 0 {
					db = db.Where("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
					continue
				}
				db = db.Or("product_id = ? and device_name = ?", d.ProductID, d.DeviceName)
			}
			return db
		}
		switch f.SharedType {
		case def.SelectTypeOnly: //直接过滤这几个设备
			db = db.WithContext(ctxs.WithAllProject(ctx)).Where(scope(db))
		case def.SelectTypeAll: //同时获取普通设备
			uc := ctxs.GetUserCtx(ctx)
			db = db.WithContext(ctxs.WithRoot(ctx)).Where("tenant_code=? and  project_id=?", uc.TenantCode, uc.ProjectID).Or(scope(db))
		}
	}
	if f.DeviceType != 0 {
		subQuery := d.db.Model(&DmProductInfo{}).Select("product_id").Where("device_type=?", f.DeviceType)
		db = db.Where("product_id in (?)", subQuery)
	}
	if f.Gateway != nil {
		subQuery := d.db.Model(&DmGatewayDevice{}).Select("product_id, device_name").
			Where(" gateway_product_id=? and gateway_device_name=?", f.Gateway.ProductID, f.Gateway.DeviceName)
		db = db.Where("(product_id, device_name) in (?)",
			subQuery)
	}
	if f.GroupID != 0 {
		subQuery := d.db.Model(&DmGroupDevice{}).Select("product_id, device_name").Where("group_id=?", f.GroupID)
		db = db.Where("(product_id, device_name) in (?)",
			subQuery)
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
func (d DeviceInfoRepo) UpdateWithField(ctx context.Context, f DeviceFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmDeviceInfo{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (d DeviceInfoRepo) UpdateOfflineStatus(ctx context.Context, f DeviceFilter) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&DmDeviceInfo{}).Updates(map[string]any{
		"is_online": def.False,
		"status":    def.DeviceStatusOffline,
	}).Error
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
