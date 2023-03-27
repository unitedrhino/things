package mysql

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/rulesvr/internal/domain/alarm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RuleAlarmRecordModel = (*customRuleAlarmRecordModel)(nil)

type (
	// RuleAlarmRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmRecordModel.
	RuleAlarmRecordModel interface {
		ruleAlarmRecordModel
		FindByFilter(ctx context.Context, filter alarm.RecordFilter, page *def.PageInfo) ([]*RuleAlarmRecord, error)
		CountByFilter(ctx context.Context, filter alarm.RecordFilter) (size int64, err error)
	}

	customRuleAlarmRecordModel struct {
		*defaultRuleAlarmRecordModel
	}
)

// NewRuleAlarmRecordModel returns a model for the database table.
func NewRuleAlarmRecordModel(conn sqlx.SqlConn) RuleAlarmRecordModel {
	return &customRuleAlarmRecordModel{
		defaultRuleAlarmRecordModel: newRuleAlarmRecordModel(conn),
	}
}

func (c customRuleAlarmRecordModel) FmtSql(sql sq.SelectBuilder, f alarm.RecordFilter) sq.SelectBuilder {
	sql = f.Time.FmtSql(sql)
	if f.AlarmID != 0 {
		sql = sql.Where("alarmID=?", f.AlarmID)
	}
	if f.TriggerType != 0 {
		sql = sql.Where("triggerType=?", f.TriggerType)
	}
	if f.ProductID != "" {
		sql = sql.Where("productID=?", f.ProductID)
	}
	if f.DeviceName != "" {
		sql = sql.Where("deviceName=?", f.DeviceName)
	}
	return sql
}

func (c customRuleAlarmRecordModel) FindByFilter(ctx context.Context, filter alarm.RecordFilter, page *def.PageInfo) ([]*RuleAlarmRecord, error) {
	var resp []*RuleAlarmRecord
	sql := sq.Select(ruleAlarmRecordRows).From(c.table).Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
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

func (c customRuleAlarmRecordModel) CountByFilter(ctx context.Context, filter alarm.RecordFilter) (size int64, err error) {
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
