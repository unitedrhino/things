package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleAlarmSceneModel = (*customRuleAlarmSceneModel)(nil)

type (
	// RuleAlarmSceneModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleAlarmSceneModel.
	RuleAlarmSceneModel interface {
		ruleAlarmSceneModel
	}

	customRuleAlarmSceneModel struct {
		*defaultRuleAlarmSceneModel
	}
)

// NewRuleAlarmSceneModel returns a model for the database table.
func NewRuleAlarmSceneModel(conn sqlx.SqlConn) RuleAlarmSceneModel {
	return &customRuleAlarmSceneModel{
		defaultRuleAlarmSceneModel: newRuleAlarmSceneModel(conn),
	}
}
