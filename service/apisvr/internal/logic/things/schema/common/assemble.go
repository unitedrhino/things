package common

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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
		Order:             in.Order,
		IsCanSceneLinkage: in.IsCanSceneLinkage,
		IsHistory:         in.IsHistory,
		IsShareAuthPerm:   in.IsShareAuthPerm,
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
		Order:             in.Order,
		IsCanSceneLinkage: in.IsCanSceneLinkage,
		IsHistory:         in.IsHistory,
		IsShareAuthPerm:   in.IsShareAuthPerm,
		Affordance:        utils.ToNullString(in.Affordance),
	}
	return &rpc
}
