package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DmProductScriptModel = (*customDmProductScriptModel)(nil)

type (
	// DmProductScriptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmProductScriptModel.
	DmProductScriptModel interface {
		dmProductScriptModel
	}

	customDmProductScriptModel struct {
		*defaultDmProductScriptModel
	}
)

// NewDmProductScriptModel returns a model for the database table.
func NewDmProductScriptModel(conn sqlx.SqlConn) DmProductScriptModel {
	return &customDmProductScriptModel{
		defaultDmProductScriptModel: newDmProductScriptModel(conn),
	}
}
