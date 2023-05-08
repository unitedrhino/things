package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CasbinRulesModel = (*customCasbinRulesModel)(nil)

type (
	// CasbinRulesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCasbinRulesModel.
	CasbinRulesModel interface {
		casbinRulesModel
	}

	customCasbinRulesModel struct {
		*defaultCasbinRulesModel
	}
)

// NewCasbinRulesModel returns a model for the database table.
func NewCasbinRulesModel(conn sqlx.SqlConn) CasbinRulesModel {
	return &customCasbinRulesModel{
		defaultCasbinRulesModel: newCasbinRulesModel(conn),
	}
}
