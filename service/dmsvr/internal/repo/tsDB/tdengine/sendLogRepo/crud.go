package sendLogRepo

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
)

func (s *SendLogRepo) fillFilter(sql sq.SelectBuilder, filter deviceLog.SendFilter) sq.SelectBuilder {
	if filter.UserID != 0 {
		sql = sql.Where("`user_id`=?", filter.UserID)
	}
	if filter.DataID != "" {
		sql = sql.Where("`data_id`=?", filter.DataID)
	}
	if len(filter.DataIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`data_id` in (%v)", stores.ArrayToSql(filter.DataIDs)))
	}
	if len(filter.ProductID) != 0 {
		sql = sql.Where("`product_id`=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		sql = sql.Where("`device_name`=?", filter.DeviceName)
	}
	if len(filter.ProductIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`product_id` in (%v)", stores.ArrayToSql(filter.ProductIDs)))
	}
	if filter.TenantCode != "" {
		sql = sql.Where("`tenant_code`=?", filter.TenantCode)
	}
	if filter.ProjectID != 0 {
		sql = sql.Where("`project_id`=?", filter.ProjectID)
	}
	if filter.AreaID != 0 {
		sql = sql.Where("`area_id`=?", filter.AreaID)
	}
	if filter.AreaIDPath != "" {
		sql = sql.Where("`area_id_path` like ?", filter.AreaIDPath+"%")
	}
	if len(filter.AreaIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`area_id` in (%v)", stores.ArrayToSql(filter.AreaIDs)))
	}
	if len(filter.Actions) != 0 {
		sql = sql.Where(fmt.Sprintf("`action` in (%v)", stores.ArrayToSql(filter.Actions)))
	}
	if filter.ResultCode != 0 {
		sql = sql.Where("`result_code`=?", filter.ResultCode)
	}
	return sql
}

func (s *SendLogRepo) GetCountLog(ctx context.Context, filter deviceLog.SendFilter, page def.PageInfo2) (int64, error) {
	sqSql := sq.Select("Count(1)").From(s.GetLogStableName())
	sqSql = s.fillFilter(sqSql, filter)
	sqSql = page.FmtWhere(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return 0, err
	}
	row := s.t.QueryRowContext(ctx, sqlStr, value...)
	if err != nil {
		return 0, err
	}
	var (
		total int64
	)

	err = row.Scan(&total)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return total, nil
}

func (s *SendLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.SendFilter, page def.PageInfo2) (
	[]*deviceLog.Send, error) {
	sql := sq.Select("*").From(s.GetLogStableName()).OrderBy("`ts` desc")
	sql = s.fillFilter(sql, filter)
	sql = page.FmtSql(sql)
	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := s.t.QueryContext(ctx, sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]any
	stores.Scan(rows, &datas)
	retLogs := make([]*deviceLog.Send, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceLog(v))
	}
	return retLogs, nil
}

func (s *SendLogRepo) Insert(ctx context.Context, data *deviceLog.Send) error {
	sql := fmt.Sprintf("  %s using %s tags('%s','%s','%s',%d,%d,'%s')(`ts`, `user_id`,`account` ,`action` ,`data_id` ,`trace_id` ,`content`,`result_code`) values (?,?,?,?,?,?,?,?) ",
		s.GetLogTableName(data.ProductID, data.DeviceName), s.GetLogStableName(), data.ProductID, data.DeviceName, data.TenantCode, data.ProjectID,
		data.AreaID, data.AreaIDPath)
	s.t.AsyncInsert(sql, data.Timestamp, data.UserID, data.Account, data.Action, data.DataID, data.TraceID, data.Content, data.ResultCode)
	return nil
}
