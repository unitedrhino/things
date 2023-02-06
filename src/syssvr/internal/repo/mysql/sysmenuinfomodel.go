package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysMenuInfoModel = (*customSysMenuInfoModel)(nil)

type (
	// SysMenuInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysMenuInfoModel.
	SysMenuInfoModel interface {
		sysMenuInfoModel
	}

	customSysMenuInfoModel struct {
		*defaultSysMenuInfoModel
	}
)

// NewSysMenuInfoModel returns a model for the database table.
func NewSysMenuInfoModel(conn sqlx.SqlConn) SysMenuInfoModel {
	return &customSysMenuInfoModel{
		defaultSysMenuInfoModel: newSysMenuInfoModel(conn),
	}
}
