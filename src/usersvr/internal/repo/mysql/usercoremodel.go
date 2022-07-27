package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserCoreModel = (*customUserCoreModel)(nil)

type (
	// UserCoreModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCoreModel.
	UserCoreModel interface {
		userCoreModel
	}

	customUserCoreModel struct {
		*defaultUserCoreModel
	}
)

// NewUserCoreModel returns a model for the database table.
func NewUserCoreModel(conn sqlx.SqlConn) UserCoreModel {
	return &customUserCoreModel{
		defaultUserCoreModel: newUserCoreModel(conn),
	}
}
