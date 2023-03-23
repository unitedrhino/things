package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleAlarmLogModel = (*customRuleAlarmLogModel)(nil)

type (
	// RuleAlarmLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmLogModel.
	RuleAlarmLogModel interface {
		ruleAlarmLogModel
		FindByFilter(ctx context.Context, filter alarm.LogFilter, page *def.PageInfo) ([]*RuleAlarmLog, error)
		CountByFilter(ctx context.Context, filter alarm.LogFilter) (size int64, err error)
	}

	customRuleAlarmLogModel struct {
		*defaultRuleAlarmLogModel
	}
)

// NewRuleAlarmLogModel returns a model for the database table.
func NewRuleAlarmLogModel(conn sqlx.SqlConn) RuleAlarmLogModel {
	return &customRuleAlarmLogModel{
		defaultRuleAlarmLogModel: newRuleAlarmLogModel(conn),
	}
}

func (c customRuleAlarmLogModel) FmtSql(sql sq.SelectBuilder, f alarm.LogFilter) sq.SelectBuilder {
	return f.Time.FmtSql(sql)
}

func (c customRuleAlarmLogModel) FindByFilter(ctx context.Context, filter alarm.LogFilter, page *def.PageInfo) ([]*RuleAlarmLog, error) {
	var resp []*RuleAlarmLog
	sql := sq.Select(ruleAlarmLogRows).From(c.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
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

func (c customRuleAlarmLogModel) CountByFilter(ctx context.Context, filter alarm.LogFilter) (size int64, err error) {
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
