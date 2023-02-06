package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleFlowInfoModel = (*customRuleFlowInfoModel)(nil)

type (
	// RuleFlowInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleFlowInfoModel.
	RuleFlowInfoModel interface {
		ruleFlowInfoModel
	}

	customRuleFlowInfoModel struct {
		*defaultRuleFlowInfoModel
	}
)

// NewRuleFlowInfoModel returns a model for the database table.
func NewRuleFlowInfoModel(conn sqlx.SqlConn) RuleFlowInfoModel {
	return &customRuleFlowInfoModel{
		defaultRuleFlowInfoModel: newRuleFlowInfoModel(conn),
	}
}
