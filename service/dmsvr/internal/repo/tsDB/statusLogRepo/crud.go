package statusLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
)

func (s *StatusLogRepo) fillFilter(ctx context.Context, db *stores.DB, filter deviceLog.StatusFilter) *stores.DB {
	db = db.WithContext(ctx)
	if len(filter.ProductID) != 0 {
		db = db.Where("product_id=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		db = db.Where("device_name=?", filter.DeviceName)
	}
	if filter.Status != 0 {
		db = db.Where("status=?", filter.Status == def.True)
	}
	return db
}

func (s *StatusLogRepo) GetCountLog(ctx context.Context, filter deviceLog.StatusFilter, page def.PageInfo2) (int64, error) {
	db := s.fillFilter(ctx, s.db, filter)
	var count int64
	err := db.Model(Status{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}

func (s *StatusLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.StatusFilter, page def.PageInfo2) (
	[]*deviceLog.Status, error) {
	db := s.fillFilter(ctx, s.db, filter)
	db = page.FmtSql2(db)
	var list []*deviceLog.Status
	err := db.Model(Status{}).Find(&list).Error
	return list, stores.ErrFmt(err)
}

func (s *StatusLogRepo) Insert(ctx context.Context, data *deviceLog.Status) error {
	s.asyncInsert.AsyncInsert(&Status{Status: data})
	return nil
}
