package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleAlarmInfoModel = (*customRuleAlarmInfoModel)(nil)

type (
	// RuleAlarmInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmInfoModel.
	RuleAlarmInfoModel interface {
		ruleAlarmInfoModel
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
