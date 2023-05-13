package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysApiAuthModel = (*customSysApiAuthModel)(nil)

type (
	// SysApiAuthModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysApiAuthModel.
	SysApiAuthModel interface {
		sysApiAuthModel
	}

	customSysApiAuthModel struct {
		*defaultSysApiAuthModel
	}
)

// NewSysApiAuthModel returns a model for the database table.
func NewSysApiAuthModel(conn sqlx.SqlConn) SysApiAuthModel {
	return &customSysApiAuthModel{
		defaultSysApiAuthModel: newSysApiAuthModel(conn),
	}
}
