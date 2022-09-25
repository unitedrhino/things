package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GroupDeviceModel = (*customGroupDeviceModel)(nil)

type (
	// GroupDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupDeviceModel.
	GroupDeviceModel interface {
		groupDeviceModel
	}

	customGroupDeviceModel struct {
		*defaultGroupDeviceModel
	}
)

// NewGroupDeviceModel returns a model for the database table.
func NewGroupDeviceModel(conn sqlx.SqlConn) GroupDeviceModel {
	return &customGroupDeviceModel{
		defaultGroupDeviceModel: newGroupDeviceModel(conn),
	}
}
