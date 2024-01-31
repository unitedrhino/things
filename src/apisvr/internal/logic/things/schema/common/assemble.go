package common

import (
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToSchemaInfoRpc(in *types.CommonSchemaInfo) *dm.CommonSchemaInfo {
	if in == nil {
		return nil
	}
	rpc := &dm.CommonSchemaInfo{
		Id:                in.ID,
		Type:              in.Type,
		Identifier:        in.Identifier,
		ExtendConfig:      in.ExtendConfig,
		Name:              utils.ToRpcNullString(in.Name),
		Desc:              utils.ToRpcNullString(in.Desc),
		Required:          in.Required,
		IsCanSceneLinkage: in.IsCanSceneLinkage,
		Affordance:        utils.ToRpcNullString(in.Affordance),
	}
	return rpc
}

func ToSchemaInfoTypes(in *dm.CommonSchemaInfo) *types.CommonSchemaInfo {
	if in == nil {
		return nil
	}
	rpc := types.CommonSchemaInfo{
		ID:                in.Id,
		Type:              in.Type,
		Identifier:        in.Identifier,
		ExtendConfig:      in.ExtendConfig,
		Name:              utils.ToNullString(in.Name),
		Desc:              utils.ToNullString(in.Desc),
		Required:          in.Required,
		IsCanSceneLinkage: in.IsCanSceneLinkage,
		Affordance:        utils.ToNullString(in.Affordance),
	}
	return &rpc
}
