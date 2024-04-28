package schemamanagelogic

import (
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToCommonSchemaRpc(info *relationDB.DmCommonSchema) *dm.CommonSchemaInfo {
	db := &dm.CommonSchemaInfo{
		Id:                info.ID,
		Type:              info.Type,
		Identifier:        info.Identifier,
		ExtendConfig:      info.ExtendConfig,
		Name:              utils.ToRpcNullString(&info.Name),
		Desc:              utils.ToRpcNullString(&info.Desc),
		Affordance:        utils.ToRpcNullString(&info.Affordance),
		Required:          info.Required,
		Order:             info.Order,
		IsCanSceneLinkage: info.IsCanSceneLinkage,
		FuncGroup:         info.FuncGroup,
		UserPerm:          info.UserPerm,
		//IsShareAuthPerm:   info.IsShareAuthPerm,
		IsHistory: info.IsHistory,
	}
	return db
}

func ToCommonSchemaPo(info *dm.CommonSchemaInfo) *relationDB.DmCommonSchema {
	db := &relationDB.DmCommonSchema{
		DmSchemaCore: relationDB.DmSchemaCore{
			Tag:               schema.TagOptional,
			Type:              info.Type,
			Identifier:        info.Identifier,
			ExtendConfig:      info.ExtendConfig,
			Name:              info.Name.GetValue(),
			Desc:              info.Desc.GetValue(),
			Required:          info.Required,
			IsCanSceneLinkage: info.IsCanSceneLinkage,
			FuncGroup:         info.FuncGroup,
			UserPerm:          info.UserPerm,
			IsHistory:         info.IsHistory,
			Order:             info.Order,
			Affordance:        info.Affordance.GetValue(),
		},
	}
	return db
}
