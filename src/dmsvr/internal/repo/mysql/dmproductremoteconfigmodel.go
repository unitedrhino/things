package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DmProductRemoteConfigModel = (*customDmProductRemoteConfigModel)(nil)

type (
	// DmProductRemoteConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmProductRemoteConfigModel.
	DmProductRemoteConfigModel interface {
		dmProductRemoteConfigModel
	}

	customDmProductRemoteConfigModel struct {
		*defaultDmProductRemoteConfigModel
	}
)

// NewDmProductRemoteConfigModel returns a model for the database table.
func NewDmProductRemoteConfigModel(conn sqlx.SqlConn) DmProductRemoteConfigModel {
	return &customDmProductRemoteConfigModel{
		defaultDmProductRemoteConfigModel: newDmProductRemoteConfigModel(conn),
	}
}
