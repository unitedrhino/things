package sdkLogRepo

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"

	"gitee.com/unitedrhino/share/def"
)

func (d SDKLogRepo) fillFilter(ctx context.Context, db *stores.DB, filter deviceLog.SDKFilter) *stores.DB {
	db = db.WithContext(ctx)
	if len(filter.ProductID) != 0 {
		db = db.Where("product_id=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		db = db.Where("device_name=?", filter.DeviceName)
	}
	if filter.LogLevel != 0 {
		db = db.Where("log_level=?", filter.LogLevel)
	}
	return db
}
func (d SDKLogRepo) GetCountLog(ctx context.Context, filter deviceLog.SDKFilter, page def.PageInfo2) (int64, error) {
	db := d.fillFilter(ctx, d.db, filter)
	var count int64
	err := db.Model(SDK{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}

func (d SDKLogRepo) GetDeviceSDKLog(ctx context.Context,
	filter deviceLog.SDKFilter, page def.PageInfo2) ([]*deviceLog.SDK, error) {
	db := d.fillFilter(ctx, d.db, filter)
	db = page.FmtSql2(db)
	var list []*deviceLog.SDK
	err := db.Model(SDK{}).Find(&list).Error
	return list, stores.ErrFmt(err)
}

func (d SDKLogRepo) Insert(ctx context.Context, data *deviceLog.SDK) error {
	d.asyncInsert.AsyncInsert(&SDK{SDK: data})
	return nil
}
