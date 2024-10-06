package hubLogRepo

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	sq "gitee.com/i-Things/squirrel"
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
)

func (h *HubLogRepo) fillFilter(sql sq.SelectBuilder, filter deviceLog.HubFilter) sq.SelectBuilder {
	if len(filter.ProductID) != 0 {
		sql = sql.Where("`product_id`=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		sql = sql.Where("`device_name`=?", filter.DeviceName)
	}
	if len(filter.Content) != 0 {
		sql = sql.Where("`content`=?", filter.Content)
	}
	if len(filter.RequestID) != 0 {
		sql = sql.Where("`request_id`=?", filter.RequestID)
	}
	if len(filter.Actions) != 0 {
		sql = sql.Where(fmt.Sprintf("`action` in (%v)", stores.ArrayToSql(filter.Actions)))
	}
	if len(filter.Topics) != 0 {
		sql = sql.Where(fmt.Sprintf("`topic` in (%v)", stores.ArrayToSql(filter.Topics)))
	}
	return sql
}

func (h *HubLogRepo) GetCountLog(ctx context.Context, filter deviceLog.HubFilter, page def.PageInfo2) (int64, error) {
	sqSql := sq.Select("Count(1)").From(h.GetLogStableName())
	sqSql = h.fillFilter(sqSql, filter)
	sqSql = page.FmtWhere(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return 0, err
	}
	row := h.t.QueryRowContext(ctx, sqlStr, value...)
	var (
		total int64
	)

	err = row.Scan(&total)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return total, nil
}

func (h *HubLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.HubFilter, page def.PageInfo2) (
	[]*deviceLog.Hub, error) {
	sqql := sq.Select("*").From(h.GetLogStableName()).OrderBy("`ts` desc")
	sqql = h.fillFilter(sqql, filter)
	sqql = page.FmtSql(sqql)
	sqlStr, value, err := sqql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := h.t.QueryContext(ctx, sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]any
	stores.Scan(rows, &datas)
	retLogs := make([]*deviceLog.Hub, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceLog(filter.ProductID, v))
	}
	return retLogs, nil
}

func (h *HubLogRepo) Insert(ctx context.Context, data *deviceLog.Hub) error {
	sql := fmt.Sprintf(" %s using %s tags('%s','%s')(`ts`, `content`, `topic`, `action`,"+
		" `request_id`, `trace_id`, `result_type`,`resp_payload`) values (?,?,?,?,?,?,?,?);",
		h.GetLogTableName(data.ProductID, data.DeviceName), h.GetLogStableName(), data.ProductID, data.DeviceName)
	//if _, err := h.t.ExecContext(ctx, sql, data.Timestamp, data.Content, data.Topic, data.Action,
	//	data.RequestID, data.TraceID, data.ResultCode); err != nil {
	//	return err
	//}
	h.t.AsyncInsert(sql, data.Timestamp, data.Content, data.Topic, data.Action,
		data.RequestID, data.TraceID, data.ResultCode, data.RespPayload)
	return nil
}
