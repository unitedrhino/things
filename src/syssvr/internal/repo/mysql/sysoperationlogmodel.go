package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysOperationLogModel = (*customSysOperationLogModel)(nil)

type (
	// SysOperationLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOperationLogModel.
	SysOperationLogModel interface {
		sysOperationLogModel
	}

	customSysOperationLogModel struct {
		*defaultSysOperationLogModel
	}
)

// NewSysOperationLogModel returns a model for the database table.
func NewSysOperationLogModel(conn sqlx.SqlConn) SysOperationLogModel {
	return &customSysOperationLogModel{
		defaultSysOperationLogModel: newSysOperationLogModel(conn),
	}
}
