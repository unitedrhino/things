package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RoleMenuModel = (*customRoleMenuModel)(nil)

type (
	// RoleMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleMenuModel.
	RoleMenuModel interface {
		roleMenuModel
	}

	customRoleMenuModel struct {
		*defaultRoleMenuModel
	}
)

// NewRoleMenuModel returns a model for the database table.
func NewRoleMenuModel(conn sqlx.SqlConn) RoleMenuModel {
	return &customRoleMenuModel{
		defaultRoleMenuModel: newRoleMenuModel(conn),
	}
}
