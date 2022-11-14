package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductRemoteConfigModel = (*customProductRemoteConfigModel)(nil)

type (
	// ProductRemoteConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductRemoteConfigModel.
	ProductRemoteConfigModel interface {
		productRemoteConfigModel
	}

	customProductRemoteConfigModel struct {
		*defaultProductRemoteConfigModel
	}
)

// NewProductRemoteConfigModel returns a model for the database table.
func NewProductRemoteConfigModel(conn sqlx.SqlConn) ProductRemoteConfigModel {
	return &customProductRemoteConfigModel{
		defaultProductRemoteConfigModel: newProductRemoteConfigModel(conn),
	}
}
