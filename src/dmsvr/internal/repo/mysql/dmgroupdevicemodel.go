package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DmGroupDeviceModel = (*customDmGroupDeviceModel)(nil)

type (
	// DmGroupDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmGroupDeviceModel.
	DmGroupDeviceModel interface {
		dmGroupDeviceModel
	}

	customDmGroupDeviceModel struct {
		*defaultDmGroupDeviceModel
	}
)

// NewDmGroupDeviceModel returns a model for the database table.
func NewDmGroupDeviceModel(conn sqlx.SqlConn) DmGroupDeviceModel {
	return &customDmGroupDeviceModel{
		defaultDmGroupDeviceModel: newDmGroupDeviceModel(conn),
	}
}
