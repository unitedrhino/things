package abnormalLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"time"
)

func (s *AbnormalLogRepo) fillFilter(ctx context.Context, db *stores.DB, filter deviceLog.AbnormalFilter) *stores.DB {
	db = db.WithContext(ctx)
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
