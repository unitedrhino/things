package common

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func ToSchemaInfoRpc(in *types.CommonSchemaInfo) *dm.CommonSchemaInfo {
	return utils.Copy[dm.CommonSchemaInfo](in)
}

func ToSchemaInfoTypes(in *dm.CommonSchemaInfo) *types.CommonSchemaInfo {
	return utils.Copy[types.CommonSchemaInfo](in)
}
