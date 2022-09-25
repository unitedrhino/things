package productmanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToProductSchemaRpc(info *mysql.ProductSchema2) *dm.ProductSchemaInfo {
	db := &dm.ProductSchemaInfo{
		ProductID:  info.ProductID,
		Tag:        info.Tag,
		Type:       info.Type,
		Identifier: info.Identifier,
		Name:       utils.ToRpcNullString(&info.Name),
		Desc:       utils.ToRpcNullString(&info.Desc),
		Required:   info.Required,
		//Affordance:  info.Affordance,
	}
	return db
}
