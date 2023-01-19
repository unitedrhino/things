package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DmGroupInfoModel = (*customDmGroupInfoModel)(nil)

type (
	// DmGroupInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDmGroupInfoModel.
	DmGroupInfoModel interface {
		dmGroupInfoModel
	}

	customDmGroupInfoModel struct {
		*defaultDmGroupInfoModel
	}
)

// NewDmGroupInfoModel returns a model for the database table.
func NewDmGroupInfoModel(conn sqlx.SqlConn) DmGroupInfoModel {
	return &customDmGroupInfoModel{
		defaultDmGroupInfoModel: newDmGroupInfoModel(conn),
	}
}
