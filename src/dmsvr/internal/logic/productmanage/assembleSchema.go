package productmanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToProductSchemaRpc(info *mysql.DmProductSchema) *dm.ProductSchemaInfo {
	db := &dm.ProductSchemaInfo{
		ProductID:  info.ProductID,
		Tag:        info.Tag,
		Type:       info.Type,
		Identifier: info.Identifier,
		Name:       utils.ToRpcNullString(&info.Name),
		Desc:       utils.ToRpcNullString(&info.Desc),
		Required:   info.Required,
		Affordance: utils.ToRpcNullString(&info.Affordance),
	}
	return db
}

func ToProductSchemaPo(info *dm.ProductSchemaInfo) *mysql.DmProductSchema {
	db := &mysql.DmProductSchema{
		ProductID:  info.ProductID,
		Tag:        info.Tag,
		Type:       info.Type,
		Identifier: info.Identifier,
		Name:       info.Name.GetValue(),
		Desc:       info.Desc.GetValue(),
		Required:   info.Required,
		Affordance: info.Affordance.GetValue(),
	}
	return db
}
