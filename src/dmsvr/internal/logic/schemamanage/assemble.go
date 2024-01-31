package schemamanagelogic

import (
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToCommonSchemaRpc(info *relationDB.DmCommonSchema) *dm.CommonSchemaInfo {
	db := &dm.CommonSchemaInfo{
		Type:         info.Type,
		Identifier:   info.Identifier,
		ExtendConfig: info.ExtendConfig,
		Name:         utils.ToRpcNullString(&info.Name),
		Desc:         utils.ToRpcNullString(&info.Desc),
		Affordance:   utils.ToRpcNullString(&info.Affordance),
		Required:     info.Required,
	}
	return db
}

func ToCommonSchemaPo(info *dm.CommonSchemaInfo) *relationDB.DmCommonSchema {
	db := &relationDB.DmCommonSchema{
		DmSchemaCore: relationDB.DmSchemaCore{
			Type:         info.Type,
			Identifier:   info.Identifier,
			ExtendConfig: info.ExtendConfig,
			Name:         info.Name.GetValue(),
			Desc:         info.Desc.GetValue(),
			Required:     info.Required,
			Affordance:   info.Affordance.GetValue(),
		},
	}
	return db
}
