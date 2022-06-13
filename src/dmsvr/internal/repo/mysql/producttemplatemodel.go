package mysql

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductTemplateModel = (*customProductTemplateModel)(nil)

type (
	// ProductTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductTemplateModel.
	ProductTemplateModel interface {
		productTemplateModel
	}

	customProductTemplateModel struct {
		*defaultProductTemplateModel
	}
)

// NewProductTemplateModel returns a model for the database table.
func NewProductTemplateModel(conn sqlx.SqlConn, c cache.CacheConf) ProductTemplateModel {
	return &customProductTemplateModel{
		defaultProductTemplateModel: newProductTemplateModel(conn, c),
	}
}
