package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ MenuInfoModel = (*customMenuInfoModel)(nil)

type (
	// MenuInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMenuInfoModel.
	MenuInfoModel interface {
		menuInfoModel
	}

	customMenuInfoModel struct {
		*defaultMenuInfoModel
	}
)

// NewMenuInfoModel returns a model for the database table.
func NewMenuInfoModel(conn sqlx.SqlConn) MenuInfoModel {
	return &customMenuInfoModel{
		defaultMenuInfoModel: newMenuInfoModel(conn),
	}
}
