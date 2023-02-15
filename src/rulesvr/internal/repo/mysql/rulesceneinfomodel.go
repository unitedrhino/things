package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleSceneInfoModel = (*customRuleSceneInfoModel)(nil)

type (
	// RuleSceneInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleSceneInfoModel.
	RuleSceneInfoModel interface {
		ruleSceneInfoModel
	}

	customRuleSceneInfoModel struct {
		*defaultRuleSceneInfoModel
	}
)

// NewRuleSceneInfoModel returns a model for the database table.
func NewRuleSceneInfoModel(conn sqlx.SqlConn) RuleSceneInfoModel {
	return &customRuleSceneInfoModel{
		defaultRuleSceneInfoModel: newRuleSceneInfoModel(conn),
	}
}
