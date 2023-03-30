package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleAlarmDealRecordModel = (*customRuleAlarmDealRecordModel)(nil)

type (
	// RuleAlarmDealRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmDealRecordModel.
	RuleAlarmDealRecordModel interface {
		ruleAlarmDealRecordModel
		FindByFilter(ctx context.Context, filter alarm.DealRecordFilter, page *def.PageInfo) ([]*RuleAlarmDealRecord, error)
		CountByFilter(ctx context.Context, filter alarm.DealRecordFilter) (size int64, err error)
	}

	customRuleAlarmDealRecordModel struct {
		*defaultRuleAlarmDealRecordModel
	}
)

// NewRuleAlarmDealRecordModel returns a model for the database table.
func NewRuleAlarmDealRecordModel(conn sqlx.SqlConn) RuleAlarmDealRecordModel {
	return &customRuleAlarmDealRecordModel{
		defaultRuleAlarmDealRecordModel: newRuleAlarmDealRecordModel(conn),
	}
}

func (c customRuleAlarmDealRecordModel) FmtSql(sql sq.SelectBuilder, f alarm.DealRecordFilter) sq.SelectBuilder {
	sql = f.Time.FmtSql(sql)
	if f.AlarmRecordID != 0 {
		sql = sql.Where("alarmRecordID=?", f.AlarmRecordID)
	}
	return sql
}

func (c customRuleAlarmDealRecordModel) FindByFilter(ctx context.Context, filter alarm.DealRecordFilter, page *def.PageInfo) ([]*RuleAlarmDealRecord, error) {
	var resp []*RuleAlarmDealRecord
	sql := sq.Select(ruleAlarmDealRecordRows).From(c.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
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

func (c customRuleAlarmDealRecordModel) CountByFilter(ctx context.Context, filter alarm.DealRecordFilter) (size int64, err error) {
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
