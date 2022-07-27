package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductSchemaModel = (*customProductSchemaModel)(nil)

type (
	// ProductSchemaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductSchemaModel.
	ProductSchemaModel interface {
		productSchemaModel
	}

	customProductSchemaModel struct {
		*defaultProductSchemaModel
	}
)

// NewProductSchemaModel returns a model for the database table.
func NewProductSchemaModel(conn sqlx.SqlConn) ProductSchemaModel {
	return &customProductSchemaModel{
		defaultProductSchemaModel: newProductSchemaModel(conn),
	}
}
