package hubLogRepo

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
)

func (h *HubLogRepo) fillFilter(ctx context.Context, db *stores.DB, filter deviceLog.HubFilter) *stores.DB {
	db = db.WithContext(ctx)
	if len(filter.ProductID) != 0 {
		db = db.Where("product_id=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		db = db.Where("device_name=?", filter.DeviceName)
	}
	if len(filter.Content) != 0 {
		db = db.Where("content=?", filter.Content)
	}
	if len(filter.RequestID) != 0 {
		db = db.Where("request_id=?", filter.RequestID)
	}
	if len(filter.Actions) != 0 {
		db = db.Where(fmt.Sprintf("action in (%v)", stores.ArrayToSql(filter.Actions)))
	}
	if len(filter.Topics) != 0 {
		db = db.Where(fmt.Sprintf("topic in (%v)", stores.ArrayToSql(filter.Topics)))
	}
	return db
}

func (h *HubLogRepo) GetCountLog(ctx context.Context, filter deviceLog.HubFilter, page def.PageInfo2) (int64, error) {
	db := h.fillFilter(ctx, h.db, filter)
	var count int64
	err := db.Model(Hub{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}

func (h *HubLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.HubFilter, page def.PageInfo2) (
	[]*deviceLog.Hub, error) {
	db := h.fillFilter(ctx, h.db, filter)
	db = page.FmtSql2(db)
	var list []*deviceLog.Hub
	err := db.Model(Hub{}).Find(&list).Error
	return list, stores.ErrFmt(err)
}

func (h *HubLogRepo) Insert(ctx context.Context, data *deviceLog.Hub) error {
	h.asyncInsert.AsyncInsert(&Hub{Hub: data})
	return nil
}
