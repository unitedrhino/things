package abnormalLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB"
	"time"
)

func (s *AbnormalLogRepo) fillFilter(ctx context.Context, db *stores.DB, filter deviceLog.AbnormalFilter) *stores.DB {
	db = db.WithContext(ctx)
	db.Statement.Dialector.Name()
	if len(filter.ProductID) != 0 {
		db = db.Where("product_id=?", filter.ProductID)
	}

	if len(filter.DeviceName) != 0 {
		db = db.Where("device_name=?", filter.DeviceName)
	}

	if filter.Action != 0 {
		db = db.Where("action=?", def.ToBool(filter.Action))
	}
	if filter.Type != "" {
		db = db.Where(fmt.Sprintf("%s=?", stores.Col("type")), filter.Type)
	}
	if filter.Reason != "" {
		db = db.Where("reason=?", filter.Reason)
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

func (s *AbnormalLogRepo) GetCountLog(ctx context.Context, filter deviceLog.AbnormalFilter, page def.PageInfo2) (int64, error) {
	db := s.fillFilter(ctx, s.db, filter)
	var count int64
	err := db.Model(Abnormal{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}

func (s *AbnormalLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.AbnormalFilter, page def.PageInfo2) (
	[]*deviceLog.Abnormal, error) {
	db := s.fillFilter(ctx, s.db, filter)
	db = page.FmtSql2(db)
	var list []*deviceLog.Abnormal
	err := db.Model(Abnormal{}).Find(&list).Error
	return list, stores.ErrFmt(err)
}

func (s *AbnormalLogRepo) Insert(ctx context.Context, data *deviceLog.Abnormal) error {
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}
	data.TraceID = utils.TraceIdFromContext(ctx)
	s.asyncInsert.AsyncInsert(&Abnormal{Abnormal: data})
	return nil
}
