package schemaDataRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
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

func (d *DeviceDataRepo) fmtSql(ctx context.Context, db *stores.DB, filter msgThing.FilterOpt) *stores.DB {
	db = d.db.WithContext(ctx)
	if filter.ProductID != "" {
		db = db.Where("product_id=? ", filter.ProductID)
	}
	if len(filter.DeviceNames) != 0 {
		db = db.Where(fmt.Sprintf("device_name= (%v)", stores.ArrayToSql(filter.DeviceNames)))
	}
	if filter.DataID != "" {
		db = db.Where("identifier=? ", filter.DataID)
	}
	if len(filter.Types) != 0 {
		db = db.Where(fmt.Sprintf("%s = (%v)", stores.Col("type"), stores.ArrayToSql(filter.Types)))
	}

	db = tsDB.GroupFilter(db, filter.BelongGroup)
	if len(filter.ProductIDs) != 0 {
		db = db.Where("product_id IN ?", filter.ProductIDs)
	} else if filter.ProductID != "" {
		db = db.Where("product_id = ?", filter.ProductID)
	}
	subQuery := d.db.Table("dm_device_info").Model(&relationDB.DmDeviceInfo{}).Select("product_id, device_name")
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
