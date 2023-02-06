package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysRoleInfoModel = (*customSysRoleInfoModel)(nil)

type (
	// SysRoleInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleInfoModel.
	SysRoleInfoModel interface {
		sysRoleInfoModel
	}

	customSysRoleInfoModel struct {
		*defaultSysRoleInfoModel
	}
)

// NewSysRoleInfoModel returns a model for the database table.
func NewSysRoleInfoModel(conn sqlx.SqlConn) SysRoleInfoModel {
	return &customSysRoleInfoModel{
		defaultSysRoleInfoModel: newSysRoleInfoModel(conn),
	}
}
