package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleAlarmDealModel = (*customRuleAlarmDealModel)(nil)

type (
	// RuleAlarmDealModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmDealModel.
	RuleAlarmDealModel interface {
		ruleAlarmDealModel
	}

	customRuleAlarmDealModel struct {
		*defaultRuleAlarmDealModel
	}
)

// NewRuleAlarmDealModel returns a model for the database table.
func NewRuleAlarmDealModel(conn sqlx.SqlConn) RuleAlarmDealModel {
	return &customRuleAlarmDealModel{
		defaultRuleAlarmDealModel: newRuleAlarmDealModel(conn),
	}
}
