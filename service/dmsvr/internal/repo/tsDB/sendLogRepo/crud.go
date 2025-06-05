package sendLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB"
)

func (s *SendLogRepo) fillFilter(ctx context.Context, db *stores.DB, filter deviceLog.SendFilter) *stores.DB {
	db = db.WithContext(ctx)
	if filter.UserID != 0 {
		db = db.Where("user_id=?", filter.UserID)
	}
	if len(filter.ProductID) != 0 {
		db = db.Where("product_id=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		db = db.Where("device_name=?", filter.DeviceName)
	}
	if len(filter.Actions) != 0 {
		db = db.Where("action in ?", filter.Actions)
	}
	db = tsDB.GroupFilter(db, filter.BelongGroup)
	if len(filter.ProductIDs) != 0 {
		db = db.Where("product_id IN ?", filter.ProductIDs)
	}
	subQuery := s.db.Table("dm_device_info").Model(&relationDB.DmDeviceInfo{}).Select("product_id, device_name")
	var hasDeviceJoin bool
	if filter.ProjectID != 0 {
		subQuery = subQuery.Where("project_id=?", filter.ProjectID)
		hasDeviceJoin = true
	}
	if filter.TenantCode != "" {
		subQuery = subQuery.Where("tenant_code=?", filter.TenantCode)
		hasDeviceJoin = true
	}
	if filter.AreaID != 0 {
		subQuery = subQuery.Where("area_id=?", filter.AreaID)
		hasDeviceJoin = true
	}
	if filter.AreaIDPath != "" {
		subQuery = subQuery.Where("area_id_path like ?", filter.AreaIDPath+"%")
		hasDeviceJoin = true
	}
	if len(filter.AreaIDs) != 0 {
		subQuery = subQuery.Where("area_id in ?", filter.AreaIDs)
		hasDeviceJoin = true
	}
	if hasDeviceJoin {
		db = db.Where("(product_id, device_name) in (?)",
			subQuery)
	}

	return db
}

func (s *SendLogRepo) GetCountLog(ctx context.Context, filter deviceLog.SendFilter, page def.PageInfo2) (int64, error) {
	db := s.fillFilter(ctx, s.db, filter)
	var count int64
	err := db.Model(Send{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}

func (s *SendLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.SendFilter, page def.PageInfo2) (
	[]*deviceLog.Send, error) {
	db := s.fillFilter(ctx, s.db, filter)
	db = page.FmtSql2(db)
	var list []*deviceLog.Send
	err := db.Model(Send{}).Find(&list).Error
	return list, stores.ErrFmt(err)
}

func (s *SendLogRepo) Insert(ctx context.Context, data *deviceLog.Send) error {
	s.asyncInsert.AsyncInsert(&Send{Send: data})
	return nil
}
