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
		TenantCode         string
		TenantCodes        []string
		ProjectIDs         []int64
		ProductID          string
		ProductIDs         []string
		AreaIDs            []int64
		NotAreaIDs         []int64
		DeviceName         string
		DeviceNames        []string
		Tags               map[string]string
		LastLoginTime      *def.TimeRange
		IsOnline           int64
		Status             int64
		Range              int64
		Position           stores.Point
		DeviceAlias        string
		Versions           []string
		Cores              []*devices.Core
		Gateway            *devices.Core
		WithProduct        bool
		ProductCategoryID  int64
		ProductCategoryIDs []int64
		SharedType         def.SelectType
		CollectType        def.SelectType
		WithManufacturer   bool
		DeviceType         int64
		DeviceTypes        []int64
		GroupID            int64
		GroupIDs           []int64
		UserID             int64
		UserIDs            []int64
		NotGroupID         int64
		NotAreaID          int64
		Distributor        *stores.IDPathFilter
		RatedPower         *stores.Cmp
		ExpTime            *stores.Cmp
		AreaIDPath         string
		HasOwner           int64 //是否被人拥有
		NotOtaJobID        int64
		NeedConfirmJobID   int64
		NeedConfirmVersion string
	}
)

func NewDeviceInfoRepo(in any) *DeviceInfoRepo {
	return &DeviceInfoRepo{db: stores.GetCommonConn(in)}
}

func (d DeviceInfoRepo) fmtFilter(ctx context.Context, f DeviceFilter) *gorm.DB {
	db := d.db.WithContext(ctx)
	uc := ctxs.GetUserCtxNoNil(ctx)
	db = f.Distributor.Filter(db, "distributor")
	db = f.RatedPower.Where(db, "rated_power")
	db = f.ExpTime.Where(db, "exp_time")
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	if f.AreaIDPath != "" {
		db = db.Where("area_id_path like ?", f.AreaIDPath+"%")
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
		db = db.Where("product_id in (?)", d.db.WithContext(ctx).Select("product_id").Model(DmProductInfo{}).Where("category_id=?", f.ProductCategoryID))
	}
	if len(f.ProductCategoryIDs) != 0 {
		db = db.Where("product_id in (?)", d.db.WithContext(ctx).Select("product_id").Model(DmProductInfo{}).Where("category_id in ?", f.ProductCategoryIDs))
	}
	//业务过滤条件
	if f.ProductID != "" {
		db = db.Where("product_id = ?", f.ProductID)
	}
	if len(f.ProductIDs) != 0 {
		db = db.Where("product_id in ?", f.ProductIDs)
	}
	if len(f.NotAreaIDs) != 0 {
		db = db.Where("area_id not in ?", f.NotAreaIDs)
	}
	if len(f.AreaIDs) != 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	if f.NotAreaID != 0 {
		db = db.Where("area_id != ?", f.NotAreaID)
	}
	if len(f.Versions) != 0 {
		db = db.Where("version in ?", f.Versions)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name like ?", "%"+f.DeviceName+"%")
	}
	//db = f.Agency.Filter("agency", db)

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

	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}

	if len(f.UserIDs) != 0 {
		db = db.Where("user_id in ?", f.UserIDs)
	}

	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	if f.HasOwner > 0 {
		switch f.HasOwner {
		case def.True:
			db = db.Where("user_id > 1")
		case def.False:
			db = db.Where("user_id <= 1")
		}
	}
	if f.NeedConfirmVersion != "" {
		db = db.Where("need_confirm_version = ?", f.NeedConfirmVersion)
	}
	if f.NeedConfirmJobID != 0 {
		db = db.Where("need_confirm_job_id = ?", f.NeedConfirmJobID)
	}

	if f.Range > 0 {
		//f.Position 形如：point(116.393 39.905)
		db = db.Where(f.Position.Range("position", f.Range))
	}

	if f.DeviceType != 0 {
		subQuery := d.db.Model(&DmProductInfo{}).Select("product_id").Where("device_type=?", f.DeviceType)
		db = db.Where("product_id in (?)", subQuery)
	}
	if len(f.DeviceTypes) != 0 {
		subQuery := d.db.Model(&DmProductInfo{}).Select("product_id").Where("device_type in ?", f.DeviceTypes)
		db = db.Where("product_id in (?)", subQuery)
	}
	if f.NotOtaJobID != 0 {
		subQuery := d.db.Model(&DmOtaFirmwareDevice{}).Select("device_name").Where("job_id = ?", f.NotOtaJobID)
		db = db.Where("device_name not in (?)", subQuery)
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
	if len(f.GroupIDs) > 0 {
		subQuery := d.db.Model(&DmGroupDevice{}).Select("product_id, device_name").Where("group_id in ?", f.GroupIDs)
		db = db.Where("(product_id, device_name) in (?)",
			subQuery)
	}
	if f.NotGroupID != 0 {
		subQuery := d.db.Model(&DmGroupDevice{}).Select("product_id, device_name").Where("group_id=?", f.NotGroupID)
		db = db.Where("(product_id, device_name) not in (?)",
			subQuery)
	}

	if f.CollectType != 0 && uc.UserID != 0 && ctxs.GetUserCtx(ctx) != nil { //如果要获取收藏的设备
		subQuery := d.db.Model(&DmUserDeviceCollect{}).Select("product_id, device_name").Where("user_id = ?", uc.UserID)
		switch f.CollectType {
		case def.SelectTypeOnly: //直接过滤这几个设备
			db = db.WithContext(ctx).Where("(product_id, device_name)  in (?)",
				subQuery)
		case def.SelectTypeAll: //同时获取普通设备
			db = db.WithContext(ctx).Or("(product_id, device_name)  in (?)",
				subQuery)
		}
	}
	if f.SharedType != 0 && uc.UserID != 0 && ctxs.GetUserCtx(ctx) != nil { //如果要获取共享设备
		subQuery := d.db.Model(&DmUserDeviceShare{}).Select("product_id, device_name").Where("shared_user_id = ?", uc.UserID)

		switch f.SharedType {
		case def.SelectTypeOnly: //直接过滤这几个设备
			db = db.WithContext(ctxs.WithAllProject(ctx)).Where("(product_id, device_name)  in (?)",
				subQuery)
		case def.SelectTypeAll: //同时获取普通设备
			pids, err := stores.GetProjectAuthIDs(ctx)
			if err != nil {
				db.AddError(err)
				return db
			}
			if len(pids) != 0 {
				db = db.WithContext(ctxs.WithAllProject(ctx)).Where("project_id in ? or (product_id, device_name)  in (?)", pids, subQuery)
			}
		}
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

func (d DeviceInfoRepo) FindByFilter(ctx context.Context, f DeviceFilter, page *stores.PageInfo) ([]*DmDeviceInfo, error) {
	var results []*DmDeviceInfo
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (d DeviceInfoRepo) FindProductIDsByFilter(ctx context.Context, f DeviceFilter) ([]string, error) {
	var results []*DmDeviceInfo
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	err := db.Distinct("product_id").Select("product_id").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return utils.ToSliceWithFunc(results, func(in *DmDeviceInfo) string {
		return in.ProductID
	}), nil
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

func (d DeviceInfoRepo) MultiUpdate(ctx context.Context, devices []*devices.Core, info *DmDeviceInfo, columns ...string) error {
	db := d.fmtFilter(ctx, DeviceFilter{Cores: devices}).Model(&DmDeviceInfo{})
	err := db.Select(columns).Updates(info).Error
	return stores.ErrFmt(err)
}
