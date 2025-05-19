package abnormalLogRepo

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/tdengine"
	"time"
)

func (s *AbnormalLogRepo) fillFilter(sql sq.SelectBuilder, filter deviceLog.AbnormalFilter) sq.SelectBuilder {
	if len(filter.ProductID) != 0 {
		sql = sql.Where("`product_id`=?", filter.ProductID)
	}
	if len(filter.ProductIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("`product_id` in (%v)", stores.ArrayToSql(filter.ProductIDs)))
	}
	if len(filter.DeviceName) != 0 {
		sql = sql.Where("`device_name`=?", filter.DeviceName)
	}
	if filter.TenantCode != "" {
		sql = sql.Where("`tenant_code`=?", filter.TenantCode)
	}
	sql = tdengine.GroupFilter(sql, s.groupConfigs, filter.BelongGroup)
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
	if filter.Action != 0 {
		sql = sql.Where("`action`=?", def.ToBool(filter.Action))
	}
	if filter.Type != "" {
		sql = sql.Where("`type`=?", filter.Type)
	}
	if filter.Reason != "" {
		sql = sql.Where("`reason`=?", filter.Reason)
	}
	return sql
}

func (s *AbnormalLogRepo) GetCountLog(ctx context.Context, filter deviceLog.AbnormalFilter, page def.PageInfo2) (int64, error) {
	sqSql := sq.Select("Count(1)").From(s.GetLogStableName())
	sqSql = s.fillFilter(sqSql, filter)
	sqSql = page.FmtWhere(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return 0, err
	}
	row := s.t.QueryRowContext(ctx, sqlStr, value...)
	var (
		total int64
	)

	err = row.Scan(&total)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return total, nil
}

func (s *AbnormalLogRepo) GetDeviceLog(ctx context.Context, filter deviceLog.AbnormalFilter, page def.PageInfo2) (
	[]*deviceLog.Abnormal, error) {
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
	retLogs := make([]*deviceLog.Abnormal, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceLog(v))
	}
	return retLogs, nil
}

func (s *AbnormalLogRepo) Insert(ctx context.Context, data *deviceLog.Abnormal) error {
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}
	data.TraceID = utils.TraceIdFromContext(ctx)
	tagKeys, tagVals := tdengine.GenTagsParams(defaultTags, s.groupConfigs, data.BelongGroup)

	sql := fmt.Sprintf("  %s using %s (%s)tags('%s','%s','%s',%d,%d,'%s' %s)(`ts`, `type`,`reason` ,`action`  ,`trace_id` ) values (?,?,?,?,?) ",
		s.GetLogTableName(data.ProductID, data.DeviceName), s.GetLogStableName(), tagKeys, data.ProductID, data.DeviceName, data.TenantCode, data.ProjectID,
		data.AreaID, data.AreaIDPath, tagVals)
	s.t.AsyncInsert(sql, data.Timestamp, data.Type, data.Reason, def.ToBool(data.Action), data.TraceID)
	return nil
}
