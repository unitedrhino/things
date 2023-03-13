package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleAlarmDealRecordModel = (*customRuleAlarmDealRecordModel)(nil)

type (
	// RuleAlarmDealRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmDealRecordModel.
	RuleAlarmDealRecordModel interface {
		ruleAlarmDealRecordModel
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
