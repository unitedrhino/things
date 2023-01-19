package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DmProductFirmwareModel = (*customDmProductFirmwareModel)(nil)

type (
	// DmProductFirmwareModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmProductFirmwareModel.
	DmProductFirmwareModel interface {
		dmProductFirmwareModel
	}

	customDmProductFirmwareModel struct {
		*defaultDmProductFirmwareModel
	}
)

// NewDmProductFirmwareModel returns a model for the database table.
func NewDmProductFirmwareModel(conn sqlx.SqlConn) DmProductFirmwareModel {
	return &customDmProductFirmwareModel{
		defaultDmProductFirmwareModel: newDmProductFirmwareModel(conn),
	}
}
