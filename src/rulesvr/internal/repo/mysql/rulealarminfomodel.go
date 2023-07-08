package mysql

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleAlarmInfoModel = (*customRuleAlarmInfoModel)(nil)

type (
	// RuleAlarmInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmInfoModel.
	RuleAlarmInfoModel interface {
		ruleAlarmInfoModel
		FindByFilter(ctx context.Context, filter alarm.InfoFilter, page *def.PageInfo) ([]*RuleAlarmInfo, error)
		CountByFilter(ctx context.Context, filter alarm.InfoFilter) (size int64, err error)
	}

	customRuleAlarmInfoModel struct {
		*defaultRuleAlarmInfoModel
	}
)

// NewRuleAlarmInfoModel returns a model for the database table.
func NewRuleAlarmInfoModel(conn sqlx.SqlConn) RuleAlarmInfoModel {
	return &customRuleAlarmInfoModel{
		defaultRuleAlarmInfoModel: newRuleAlarmInfoModel(conn),
	}
}
func (c *customRuleAlarmInfoModel) FmtSql(sql sq.SelectBuilder, f alarm.InfoFilter) sq.SelectBuilder {
	if f.Name != "" {
		sql = sql.Where("name=?", f.Name)
	}
	if len(f.AlarmIDs) != 0 {
		sql = sql.Where(fmt.Sprintf("id in (%v)", stores.ArrayToSql(f.AlarmIDs)))
	}
	if f.SceneID != 0 {
		sql = sql.LeftJoin(fmt.Sprintf("`rule_alarm_scene` as ras on ras.alarmID=%s.id", c.table))
		sql = sql.Where("ras.sceneID=?", f.SceneID)
	}
	return sql
}

func (c *customRuleAlarmInfoModel) FindByFilter(
	ctx context.Context, filter alarm.InfoFilter, page *def.PageInfo) ([]*RuleAlarmInfo, error) {
	var resp []*RuleAlarmInfo
	sql := sq.Select(c.table + ".*").From(c.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	sql = c.FmtSql(sql, filter)
	query, arg, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	err = c.conn.QueryRowsCtx(ctx, &resp, query, arg...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
func (c *customRuleAlarmInfoModel) CountByFilter(ctx context.Context, filter alarm.InfoFilter) (size int64, err error) {
	sql := sq.Select("count(1)").From(c.table)
	sql = c.FmtSql(sql, filter)
	query, arg, err := sql.ToSql()
	if err != nil {
		return 0, err
	}
	err = c.conn.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	default:
		return 0, err
	}
}
