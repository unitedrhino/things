package mysql

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductFirmwareModel = (*customProductFirmwareModel)(nil)

type (
	// ProductFirmwareModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductFirmwareModel.
	ProductFirmwareModel interface {
		productFirmwareModel
	}

	customProductFirmwareModel struct {
		*defaultProductFirmwareModel
	}
)

// NewProductFirmwareModel returns a model for the database table.
func NewProductFirmwareModel(conn sqlx.SqlConn, c cache.CacheConf) ProductFirmwareModel {
	return &customProductFirmwareModel{
		defaultProductFirmwareModel: newProductFirmwareModel(conn, c),
	}
}
