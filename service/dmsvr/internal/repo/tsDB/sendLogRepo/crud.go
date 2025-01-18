package sendLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
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
