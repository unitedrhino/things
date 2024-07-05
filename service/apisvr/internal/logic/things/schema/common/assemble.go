package common

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
)

func ToSchemaInfoRpc(in *types.CommonSchemaInfo) *dm.CommonSchemaInfo {
	return utils.Copy[dm.CommonSchemaInfo](in)
}

func ToSchemaInfoTypes(in *dm.CommonSchemaInfo) *types.CommonSchemaInfo {
	return utils.Copy[types.CommonSchemaInfo](in)
}
