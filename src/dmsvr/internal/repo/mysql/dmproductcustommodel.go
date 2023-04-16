package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DmProductCustomModel = (*customDmProductCustomModel)(nil)

type (
	// DmProductCustomModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmProductCustomModel.
	DmProductCustomModel interface {
		dmProductCustomModel
	}

	customDmProductCustomModel struct {
		*defaultDmProductCustomModel
	}
)

// NewDmProductCustomModel returns a model for the database table.
func NewDmProductCustomModel(conn sqlx.SqlConn) DmProductCustomModel {
	return &customDmProductCustomModel{
		defaultDmProductCustomModel: newDmProductCustomModel(conn),
	}
}
