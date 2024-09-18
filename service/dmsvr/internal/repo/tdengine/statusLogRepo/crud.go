package statusLogRepo

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	sq "gitee.com/i-Things/squirrel"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"time"
)

func (s *StatusLogRepo) fillFilter(sql sq.SelectBuilder, filter deviceLog.StatusFilter) sq.SelectBuilder {
	if len(filter.ProductID) != 0 {
		sql = sql.Where("`product_id`=?", filter.ProductID)
	}
	if len(filter.DeviceName) != 0 {
		sql = sql.Where("`device_name`=?", filter.DeviceName)
	}
	if filter.Status != 0 {
		sql = sql.Where("`status`=?", filter.Status == def.True)
	}
	return sql
}

func (s *StatusLogRepo) GetCountLog(ctx context.Context, filter deviceLog.StatusFilter, page def.PageInfo2) (int64, error) {
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

func (s *StatusLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.StatusFilter, page def.PageInfo2) (
	[]*deviceLog.Status, error) {
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
	retLogs := make([]*deviceLog.Status, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceLog(v))
	}
	return retLogs, nil
}

func (s *StatusLogRepo) Insert(ctx context.Context, data *deviceLog.Status) error {
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}
	sql := fmt.Sprintf("  %s using %s tags('%s','%s') (`ts`, `status`) values (?,?) ",
		s.GetLogTableName(data.ProductID, data.DeviceName), s.GetLogStableName(), data.ProductID, data.DeviceName)
	s.t.AsyncInsert(sql, data.Timestamp, data.Status)
	return nil
}
