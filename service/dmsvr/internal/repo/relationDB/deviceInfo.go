package relationDB

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
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
		AreaID             int64
		AreaIDs            []int64
		AreaIDPath         string
		AreaIDPaths        []string
		NotAreaIDs         []int64
		DeviceName         string
		DeviceNames        []string
		Tags               map[string]string
		TagsLike           map[string]string
		LastLoginTime      *def.TimeRange
		IsOnline           int64
		Status             def.DeviceStatus
		Statuses           []def.DeviceStatus
		Range              int64
		Position           stores.Point
		DeviceAlias        string
		Versions           []string
		NotVersion         string
		Device             *devices.Core
		Cores              []*devices.Core
		Gateway            *devices.Core
		Iccid              string
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
		GroupIDPath        string
		GroupIDPaths       []string
		ParentGroupID      int64
		GroupName          string
		UserID             int64
		UserIDs            []int64
		NotGroupID         int64
		NotAreaID          int64
		Distributor        *stores.IDPathFilter
		Property           map[string]string
		RatedPower         *stores.Cmp
		ExpTime            *stores.Cmp
		Rssi               *stores.Cmp
		HasOwner           int64 //是否被人拥有
		NeedConfirmJobID   int64
		NeedConfirmVersion string
		NetType            int64
		ProtocolCode       string
		tableAlias         string
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
	db = f.Rssi.Where(db, "rssi")
	if f.WithProduct {
		db = db.Preload("ProductInfo")
	}
	if f.Iccid != "" {
		db = db.Where("iccid = ?", f.Iccid)
	}
	if f.AreaIDPath != "" {
		db = db.Where("area_id_path like ?", f.AreaIDPath+"%")
	}
	if len(f.AreaIDPaths) != 0 {
		or := db
		for _, v := range f.AreaIDPaths {
			or = or.Or("area_id_path like ?", v+"%")
		}
		db = db.Where(or)
	}
	if f.WithManufacturer {
		db = db.Preload("Manufacturer")
	}
	if len(f.TenantCodes) != 0 {
		db = db.Where("tenant_code in ?", f.TenantCodes)
	}
	if len(f.ProjectIDs) != 0 {
		db = db.Where("project_id IN ?", f.ProjectIDs)
	}
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	var productSelect = d.db.WithContext(ctx).Select("product_id").Model(DmProductInfo{})
	var hasProductFilter bool
	if f.ProductCategoryID != 0 {
		hasProductFilter = true
		productSelect = productSelect.Where("category_id=?", f.ProductCategoryID)
	}
	if len(f.ProductCategoryIDs) != 0 {
		hasProductFilter = true
		productSelect = productSelect.Where("category_id in ?", f.ProductCategoryIDs)
	}
	if f.NetType != 0 {
		productSelect = productSelect.Where("net_type = ?", f.NetType)
	}
	if hasProductFilter {
		db = db.Where("product_id in (?)", productSelect)
	}
	//业务过滤条件
	if f.ProductID != "" {
		db = db.Where(fmt.Sprintf("%s = ?", stores.ColWithT("product_id", f.tableAlias)), f.ProductID)
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
	if f.AreaID != 0 {
		db = db.Where("area_id = ?", f.AreaID)
	}
	if f.NotAreaID != 0 {
		db = db.Where("area_id != ?", f.NotAreaID)
	}
	if len(f.Versions) != 0 {
		db = db.Where(fmt.Sprintf("%s in ?", stores.ColWithT("version", f.tableAlias)), f.Versions)
	}
	if f.NotVersion != "" {
		db = db.Where("version != ?", f.NotVersion)
	}
	if f.DeviceName != "" {
		db = db.Where("device_name like ?", "%"+f.DeviceName+"%")
	}
	//db = f.Agency.Filter("agency", db)
	if f.Device != nil {
		db = db.Where("product_id = ? and device_name = ?", f.Device.ProductID, f.Device.DeviceName)
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
			db = stores.CmpJsonObjEq(k, v).Where(db, "tags")
		}
	}
	if f.TagsLike != nil {
		for k, v := range f.TagsLike {
			db = stores.CmpJsonObjLike(k, v).Where(db, "tags")
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
	if len(f.Statuses) != 0 {
		db = db.Where("status in ?", f.Statuses)
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
	var productQuery = d.db.Model(&DmProductInfo{}).Select("product_id")
	var hasProductQuery bool
	if f.DeviceType != 0 {
		productQuery = productQuery.Where("device_type=?", f.DeviceType)
		hasProductQuery = true
	}
	if f.ProtocolCode != "" {
		productQuery = productQuery.Where("protocol_code=? or sub_protocol_code=?", f.ProtocolCode, f.ProtocolCode)
		hasProductQuery = true
	}
	if len(f.DeviceTypes) != 0 {
		productQuery = productQuery.Where("device_type in ?", f.DeviceTypes)
		hasProductQuery = true
	}
	if hasProductQuery {
		db = db.Where("product_id in (?)", productQuery)
	}
	if len(f.Property) != 0 {
		subQuery := d.db.Model(&DmDeviceShadow{}).Select("product_id, device_name")
		for k, v := range f.Property {
			subQuery = subQuery.Where(" data_id=? and value=?", k, v)
		}
		db = db.Where("(product_id, device_name) in (?)",
			subQuery)
	}
	if f.Gateway != nil {
		subQuery := d.db.Model(&DmGatewayDevice{}).Select("product_id, device_name").
			Where(" gateway_product_id=? and gateway_device_name=?", f.Gateway.ProductID, f.Gateway.DeviceName)
		db = db.Where("(product_id, device_name) in (?)",
			subQuery)
	}

	var groupSubQuery = d.db.Model(&DmGroupDevice{}).Select("product_id, device_name")
	var groupFilter bool
	if f.GroupID != 0 {
		groupFilter = true
		groupSubQuery = groupSubQuery.Where("group_id=?", f.GroupID)
	}
	if len(f.GroupIDs) > 0 {
		groupFilter = true
		groupSubQuery = groupSubQuery.Where("group_id in ?", f.GroupIDs)
	}
	if f.GroupIDPath != "" {
		groupFilter = true
		groupSubQuery = groupSubQuery.Where("group_id_path like ?", f.GroupIDPath+"%")
	}
	if len(f.GroupIDPaths) != 0 {
		groupFilter = true
		or := groupSubQuery
		for _, v := range f.GroupIDPaths {
			or = or.Or("group_id_path like ?", v+"%")
		}
		groupSubQuery = groupSubQuery.Where(or)
	}

	if f.ParentGroupID != 0 || f.GroupName != "" {
		groupFilter = true
		subQuery := d.db.Model(&DmGroupInfo{}).Select("id")
		if f.ParentGroupID != 0 {
			subQuery = subQuery.Where("parent_id = ?", f.ParentGroupID)
		}
		if f.GroupName != "" {
			subQuery = subQuery.Where("name like ?", "%"+f.GroupName+"%")
		}
		groupSubQuery = groupSubQuery.Where("group_id in (?)", subQuery)
	}
	if groupFilter {
		db = db.Where("(product_id, device_name)  in (?)",
			groupSubQuery)
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
			db = db.WithContext(ctxs.WithAllProject(ctxs.WithAllArea(ctx))).Where("(product_id, device_name)  in (?)",
				subQuery)
		case def.SelectTypeAll: //同时获取普通设备
			if !(uc.IsAdmin && (uc.ProjectID <= def.NotClassified || uc.AllProject)) {
				d := shareAll(ctx, d.db, uc)
				db = db.WithContext(ctxs.WithAllProject(ctxs.WithAllArea(ctx))).Where(d)
				//if uc.ProjectID <= def.NotClassified { //如果不是管理员又没有传项目ID,则只获取分享的设备
				//	db = db.WithContext(ctxs.WithAllProject(ctxs.WithAllArea(ctx))).Where("(product_id, device_name)  in (?)",
				//		subQuery)
				//} else { //如果传了项目ID,则判断项目的权限
				//	or := d.db
				//	or = or.Or("(product_id, device_name)  in (?)", subQuery)
				//	pa := uc.ProjectAuth[uc.ProjectID]
				//	if pa == nil && uc.IsAdmin {
				//		pa = &ctxs.ProjectAuth{AuthType: def.AuthAdmin}
				//	}
				//	if pa == nil {
				//		db = db.WithContext(ctxs.WithAllProject(ctxs.WithAllArea(ctx))).Where("(product_id, device_name)  in (?)",
				//			subQuery)
				//	} else {
				//		if pa.AuthType < def.AuthRead || uc.AllArea {
				//			or = or.Or("project_id = ?", uc.ProjectID)
				//			db = db.WithContext(ctxs.WithAllProject(ctxs.WithAllArea(ctx))).Where(or)
				//		} else { //如果是读权限,还需要过滤区域
				//			areaIDs, er1 := stores.GetAreaAuthIDs(ctx)
				//			areaIDPaths, er2 := stores.GetAreaAuthIDPaths(ctx)
				//			if (er1 == nil || er2 == nil) && (areaIDPaths != nil || areaIDs != nil) {
				//				or = or.Or("(area_id in ? or area_id_path in ?) and project_id = ?", areaIDs, areaIDPaths, uc.ProjectID)
				//			}
				//			db = db.WithContext(ctxs.WithAllProject(ctxs.WithAllArea(ctx))).Where(or)
				//		}
				//	}

				//}
			}
		}
	}

	return db
}

func shareAll(ctx context.Context, db *gorm.DB, uc *ctxs.UserCtx) *gorm.DB {
	subQuery := db.Model(&DmUserDeviceShare{}).Select("product_id, device_name").Where("shared_user_id = ?", uc.UserID)
	var pas = uc.ProjectAuth
	if uc.ProjectID > def.NotClassified {
		pas = map[int64]*ctxs.ProjectAuth{uc.ProjectID: uc.ProjectAuth[uc.ProjectID]}
	}
	or := db
	or = or.Or("(product_id, device_name)  in (?)", subQuery)
	for pid, pa := range pas {
		if pa == nil && uc.IsAdmin {
			pa = &ctxs.ProjectAuth{AuthType: def.AuthAdmin}
		}
		if pa == nil {
			db = db.Where("(product_id, device_name)  in (?)",
				subQuery)
			return db
		} else {
			if pa.AuthType < def.AuthRead || uc.AllArea {
				or = or.Or("project_id = ?", pid)
			} else { //如果是读权限,还需要过滤区域
				areaIDs := utils.SetToSlice(pa.Area)
				areaIDPaths := utils.SetToSlice(pa.AreaPath)
				if areaIDPaths != nil || areaIDs != nil {
					or = or.Or("(area_id in ? or area_id_path in ?) and project_id = ?", areaIDs, areaIDPaths, pid)
				}
			}
		}
	}
	return or
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
	db := d.db.WithContext(ctx)
	if !ctxs.GetUserCtxOrNil(ctx).AllTenant {
		db = db.Omit("tenant_code")
	}
	err := db.Where("id = ?", data.ID).Save(data).Error
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

func (d DeviceInfoRepo) FindWithNotOtaJobIDByFilter(ctx context.Context, f DeviceFilter, notOtaJobID int64, page *stores.PageInfo) ([]*DmDeviceInfo, error) {
	var results []*DmDeviceInfo
	f.tableAlias = "di"
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{}).Table("dm_device_info di").
		Joins("LEFT JOIN dm_ota_firmware_device ofd ON di.device_name = ofd.device_name   AND ofd.job_id = ?  AND ofd.deleted_time = 0", notOtaJobID).
		Where("ofd.device_name IS NULL")
	db = page.ToGorm(db)
	err := db.Select("di.*").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (d DeviceInfoRepo) FindCoreByFilter(ctx context.Context, f DeviceFilter, page *stores.PageInfo) ([]devices.Core, error) {
	var results []*DmDeviceInfo
	db := d.fmtFilter(ctx, f).Model(&DmDeviceInfo{})
	db = page.ToGorm(db)
	err := db.Select("product_id", "device_name").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return utils.ToSliceWithFunc(results, func(in *DmDeviceInfo) devices.Core {
		return devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}
	}), nil
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
