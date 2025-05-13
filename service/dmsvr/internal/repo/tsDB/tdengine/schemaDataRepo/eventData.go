package schemaDataRepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
)

func (d *DeviceDataRepo) InsertEventData(ctx context.Context, productID string,
	deviceName string, event *msgThing.EventData) error {
	param, err := json.Marshal(event.Params)
	if err != nil {
		return errors.System.AddDetail("param json parse failure")
	}
	sql := fmt.Sprintf(
		" %s using %s tags('%s','%s','%s',%d,%d,'%s','%s','%s') (`ts`,`event_id`,`event_type`, `param`) values (?,?,?,?);",
		d.GetEventTableName(productID, deviceName), d.GetEventStableName(), productID, deviceName,
		event.TenantCode, event.ProjectID, event.AreaID, event.AreaIDPath, utils.GenSliceStr(event.GroupIDs),
		utils.GenSliceStr(event.GroupIDPaths))
	//if _, err := d.t.ExecContext(ctx, sql, event.TimeStamp, event.Identifier, event.Type, param); err != nil {
	//	return err
	//}
	d.t.AsyncInsert(sql, event.TimeStamp, event.Identifier, event.Type, param)
	return nil
}

func (d *DeviceDataRepo) fmtSql(f msgThing.FilterOpt, sql sq.SelectBuilder) sq.SelectBuilder {
	if f.ProductID != "" {
		sql = sql.Where("`product_id`=? ", f.ProductID)
	}
	if len(f.DeviceNames) != 0 {
		sql = sql.Where(fmt.Sprintf("`device_name` in (%v)", stores.ArrayToSql(f.DeviceNames)))
	}
	if f.DataID != "" {
		sql = sql.Where("`event_id`=? ", f.DataID)
	}
	if len(f.Types) != 0 {
		sql = sql.Where(fmt.Sprintf("`event_type` = (%v)", stores.ArrayToSql(f.Types)))
	}

	if f.TenantCode != "" {
		sql = sql.Where("`tenant_code`=?", f.TenantCode)
	}
	if len(f.GroupIDs) != 0 {
		sql = sql.Where(stores.ArrayEqToSql("group_ids", f.GroupIDs))
	}
	if len(f.GroupIDPaths) != 0 {
		sql = sql.Where(stores.ArrayEqToSql("group_ids", f.GroupIDs))
	}
	if f.ProjectID != 0 {
		sql = sql.Where("`project_id`=?", f.ProjectID)
	}
	if f.AreaID != 0 {
		sql = sql.Where("`area_id`=?", f.AreaID)
	}
	if f.AreaIDPath != "" {
		sql = sql.Where("`area_id_path` like ?", f.AreaIDPath+"%")
	}
	if len(f.AreaIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`area_id` in (%v)", stores.ArrayToSql(f.AreaIDs)))
	}
	return sql
}

func (d *DeviceDataRepo) GetEventDataByFilter(
	ctx context.Context,
	filter msgThing.FilterOpt) ([]*msgThing.EventData, error) {
	sql := sq.Select("*").From(d.GetEventStableName()).OrderBy("`ts` desc")
	sql = d.fmtSql(filter, sql)
	sql = filter.Page.FmtSql(sql)
	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.QueryContext(ctx, sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]any
	stores.Scan(rows, &datas)
	retEvents := make([]*msgThing.EventData, 0, len(datas))
	for _, v := range datas {
		retEvents = append(retEvents, ToEventData(v))
	}
	return retEvents, nil
}

func (d *DeviceDataRepo) GetEventCountByFilter(
	ctx context.Context,
	filter msgThing.FilterOpt) (int64, error) {
	sqSql := sq.Select("count(1)").From(d.GetEventStableName())
	sqSql = d.fmtSql(filter, sqSql)
	sqSql = filter.Page.FmtWhere(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return 0, err
	}
	row := d.t.QueryRowContext(ctx, sqlStr, value...)
	var total int64
	err = row.Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
