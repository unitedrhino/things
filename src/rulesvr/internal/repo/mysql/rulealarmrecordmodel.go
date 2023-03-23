package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleAlarmRecordModel = (*customRuleAlarmRecordModel)(nil)

type (
	// RuleAlarmRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmRecordModel.
	RuleAlarmRecordModel interface {
		ruleAlarmRecordModel
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
