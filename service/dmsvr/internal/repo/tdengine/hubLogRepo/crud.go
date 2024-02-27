package hubLogRepo

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/deviceMsg/msgHubLog"
	"gitee.com/i-Things/share/stores"
	sq "github.com/Masterminds/squirrel"
)

func (d HubLogRepo) fillFilter(sql sq.SelectBuilder, filter msgHubLog.HubFilter) sq.SelectBuilder {
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

func (d HubLogRepo) GetCountLog(ctx context.Context, filter msgHubLog.HubFilter, page def.PageInfo2) (int64, error) {
	sqSql := sq.Select("Count(1)").From(d.GetLogStableName())
	sqSql = d.fillFilter(sqSql, filter)
	sqSql = page.FmtWhere(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return 0, err
	}
	row := d.t.QueryRowContext(ctx, sqlStr, value...)
	if err != nil {
		return 0, err
	}
	var (
		total int64
	)

	err = row.Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func (d HubLogRepo) GetDeviceLog(ctx context.Context, filter msgHubLog.HubFilter, page def.PageInfo2) (
	[]*msgHubLog.HubLog, error) {
	sql := sq.Select("*").From(d.GetLogStableName()).OrderBy("`ts` desc")
	sql = d.fillFilter(sql, filter)
	sql = page.FmtSql(sql)
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
	retLogs := make([]*msgHubLog.HubLog, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceLog(filter.ProductID, v))
	}
	return retLogs, nil
}

func (d HubLogRepo) Insert(ctx context.Context, data *msgHubLog.HubLog) error {
	sql := fmt.Sprintf(" %s using %s tags('%s','%s')(`ts`, `content`, `topic`, `action`,"+
		" `request_id`, `trance_id`, `result_type`) values (?,?,?,?,?,?,?);",
		d.GetLogTableName(data.ProductID, data.DeviceName), d.GetLogStableName(), data.ProductID, data.DeviceName)
	//if _, err := d.t.ExecContext(ctx, sql, data.Timestamp, data.Content, data.Topic, data.Action,
	//	data.RequestID, data.TranceID, data.ResultType); err != nil {
	//	return err
	//}
	d.t.AsyncInsert(sql, data.Timestamp, data.Content, data.Topic, data.Action,
		data.RequestID, data.TranceID, data.ResultType)
	return nil
}
