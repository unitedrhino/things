package schemaDataRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
)

func (d *DeviceDataRepo) InsertEventData(ctx context.Context, productID string,
	deviceName string, event *msgThing.EventData) error {
	param, err := json.Marshal(event.Params)
	if err != nil {
		return errors.System.AddDetail("param json parse failure")
	}
	d.asyncEventInsert.AsyncInsert(&Event{
		ProductID:  productID,
		DeviceName: deviceName,
		Identifier: event.Identifier,
		Type:       event.Type,
		Param:      string(param),
		Timestamp:  event.TimeStamp,
	})
	return nil
}

func (d *DeviceDataRepo) fmtSql(ctx context.Context, db *stores.DB, f msgThing.FilterOpt) *stores.DB {
	db = d.db.WithContext(ctx)
	if f.ProductID != "" {
		db = db.Where("product_id=? ", f.ProductID)
	}
	if len(f.DeviceNames) != 0 {
		db = db.Where(fmt.Sprintf("device_name= (%v)", stores.ArrayToSql(f.DeviceNames)))
	}
	if f.DataID != "" {
		db = db.Where("identifier=? ", f.DataID)
	}
	if len(f.Types) != 0 {
		db = db.Where(fmt.Sprintf("%s = (%v)", stores.Col("type"), stores.ArrayToSql(f.Types)))
	}
	return db
}

func (d *DeviceDataRepo) GetEventDataByFilter(
	ctx context.Context,
	filter msgThing.FilterOpt) ([]*msgThing.EventData, error) {
	db := d.fmtSql(ctx, d.db, filter)
	db = filter.Page.FmtSql2(db)
	var list []*msgThing.EventData
	err := db.Model(Event{}).Find(&list).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return list, stores.ErrFmt(err)
}

func (d *DeviceDataRepo) GetEventCountByFilter(
	ctx context.Context,
	filter msgThing.FilterOpt) (int64, error) {
	db := d.fmtSql(ctx, d.db, filter)
	var count int64
	err := db.Model(Event{}).Count(&count).Error
	return count, stores.ErrFmt(err)
}
