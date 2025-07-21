package schemamanagelogic

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/schema"
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
		ControlMode:       info.ControlMode,
		UserPerm:          info.UserPerm,
		//IsShareAuthPerm:   info.IsShareAuthPerm,
		IsPassword: info.IsPassword,
		IsHistory:  info.IsHistory,
	}
	return db
}

func ToCommonSchemaPo(info *dm.CommonSchemaInfo) *relationDB.DmCommonSchema {
	db := &relationDB.DmCommonSchema{
		Identifier: info.Identifier,
		DmSchemaCore: relationDB.DmSchemaCore{
			Tag:               schema.TagOptional,
			Type:              info.Type,
			ExtendConfig:      info.ExtendConfig,
			Name:              info.Name.GetValue(),
			Desc:              info.Desc.GetValue(),
			Required:          info.Required,
			IsCanSceneLinkage: info.IsCanSceneLinkage,
			FuncGroup:         info.FuncGroup,
			ControlMode:       info.ControlMode,
			UserPerm:          info.UserPerm,
			IsHistory:         info.IsHistory,
			IsPassword:        info.IsPassword,
			Order:             info.Order,
			Affordance:        info.Affordance.GetValue(),
		},
	}
	return db
}
