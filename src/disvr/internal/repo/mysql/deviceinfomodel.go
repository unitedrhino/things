package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DeviceInfoModel = (*customDeviceInfoModel)(nil)

type (
	// DeviceInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDeviceInfoModel.
	DeviceInfoModel interface {
		deviceInfoModel
	}

	customDeviceInfoModel struct {
		*defaultDeviceInfoModel
	}
)

// NewDeviceInfoModel returns a model for the database table.
func NewDeviceInfoModel(conn sqlx.SqlConn) DeviceInfoModel {
	return &customDeviceInfoModel{
		defaultDeviceInfoModel: newDeviceInfoModel(conn),
	}
}
