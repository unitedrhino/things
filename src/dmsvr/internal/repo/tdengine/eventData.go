package tdengine

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
)

func (d *DeviceDataRepo) InsertEventData(ctx context.Context, productID string,
	deviceName string, event *deviceTemplate.EventData) error {
	param, err := json.Marshal(event.Params)
	if err != nil {
		return errors.System.AddDetail("param json parse failure")
	}
	sql := fmt.Sprintf("insert into %s (`ts`,`event_id`,`event_type`, `param`) values (?,?,?,?);", getEventTableName(productID, deviceName))
	if _, err := d.t.Exec(sql, event.TimeStamp, event.ID, event.Type, param); err != nil {
		return err
	}
	return nil
}

func (d *DeviceDataRepo) GetEventDataByID(
	ctx context.Context,
	productID string,
	deviceName string,
	dataID string,
	page def.PageInfo2) ([]*deviceTemplate.EventData, error) {
	sql := sq.Select("*").From(getEventStableName(productID)).
		Where("`device_name`=? and `event_id`=? ", deviceName, dataID).OrderBy("`ts` desc")
	sql = page.FmtSql(sql)
	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.Query(sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]interface{}
	store.Scan(rows, &datas)
	retEvents := make([]*deviceTemplate.EventData, 0, len(datas))
	for _, v := range datas {
		retEvents = append(retEvents, ToEventData(dataID, v))
	}
	return retEvents, nil
}
