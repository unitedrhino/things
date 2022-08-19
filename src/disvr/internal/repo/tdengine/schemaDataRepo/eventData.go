package schemaDataRepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
)

func (d *SchemaDataRepo) InsertEventData(ctx context.Context, productID string,
	deviceName string, event *deviceMsg.EventData) error {
	param, err := json.Marshal(event.Params)
	if err != nil {
		return errors.System.AddDetail("param json parse failure")
	}
	sql := fmt.Sprintf("insert into %s (`ts`,`event_id`,`event_type`, `param`) values (?,?,?,?);",
		d.GetEventTableName(productID, deviceName))
	if _, err := d.t.ExecContext(ctx,sql, event.TimeStamp, event.ID, event.Type, param); err != nil {
		return err
	}
	return nil
}

func (d *SchemaDataRepo) GetEventDataByID(
	ctx context.Context,
	filter deviceMsg.FilterOpt) ([]*deviceMsg.EventData, error) {
	sql := sq.Select("*").From(d.GetEventStableName(filter.ProductID)).
		Where("`device_name`=? and `event_id`=? ", filter.DeviceName, filter.DataID).OrderBy("`ts` desc")
	sql = filter.Page.FmtSql(sql)
	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.QueryContext(ctx,sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]any
	store.Scan(rows, &datas)
	retEvents := make([]*deviceMsg.EventData, 0, len(datas))
	for _, v := range datas {
		retEvents = append(retEvents, ToEventData(filter.DataID, v))
	}
	return retEvents, nil
}

func (d *SchemaDataRepo) GetEventCountByID(
	ctx context.Context,
	filter deviceMsg.FilterOpt) (int64, error) {
	sqSql := sq.Select("count(1)").From(d.GetEventStableName(filter.ProductID)).
		Where("`device_name`=? and `event_id`=? ", filter.DeviceName, filter.DataID)
	sqSql = filter.Page.FmtWhere(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return 0, err
	}
	row := d.t.QueryRowContext(ctx,sqlStr, value...)
	var total int64
	err = row.Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
