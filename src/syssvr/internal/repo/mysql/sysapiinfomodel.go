package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysApiInfoModel = (*customSysApiInfoModel)(nil)

type (
	// SysApiInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysApiInfoModel.
	SysApiInfoModel interface {
		sysApiInfoModel
	}

	customSysApiInfoModel struct {
		*defaultSysApiInfoModel
	}
)

// NewSysApiInfoModel returns a model for the database table.
func NewSysApiInfoModel(conn sqlx.SqlConn) SysApiInfoModel {
	return &customSysApiInfoModel{
		defaultSysApiInfoModel: newSysApiInfoModel(conn),
	}
}
