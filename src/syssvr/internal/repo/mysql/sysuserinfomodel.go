package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysUserInfoModel = (*customSysUserInfoModel)(nil)

type (
	// SysUserInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserInfoModel.
	SysUserInfoModel interface {
		sysUserInfoModel
	}

	customSysUserInfoModel struct {
		*defaultSysUserInfoModel
	}
)

// NewSysUserInfoModel returns a model for the database table.
func NewSysUserInfoModel(conn sqlx.SqlConn) SysUserInfoModel {
	return &customSysUserInfoModel{
		defaultSysUserInfoModel: newSysUserInfoModel(conn),
	}
}
