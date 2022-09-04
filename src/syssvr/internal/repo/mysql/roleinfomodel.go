package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RoleInfoModel = (*customRoleInfoModel)(nil)

type (
	// RoleInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleInfoModel.
	RoleInfoModel interface {
		roleInfoModel
	}

	customRoleInfoModel struct {
		*defaultRoleInfoModel
	}
)

// NewRoleInfoModel returns a model for the database table.
func NewRoleInfoModel(conn sqlx.SqlConn) RoleInfoModel {
	return &customRoleInfoModel{
		defaultRoleInfoModel: newRoleInfoModel(conn),
	}
}
