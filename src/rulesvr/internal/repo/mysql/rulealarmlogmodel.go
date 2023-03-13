package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleAlarmLogModel = (*customRuleAlarmLogModel)(nil)

type (
	// RuleAlarmLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmLogModel.
	RuleAlarmLogModel interface {
		ruleAlarmLogModel
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
