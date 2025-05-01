package info

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func ToGroupInfosTypes(in []*dm.GroupInfo) []*types.GroupInfo {
	return utils.CopySlice[types.GroupInfo](in)

}

func ToGroupInfoTypes(in *dm.GroupInfo) *types.GroupInfo {
	return utils.Copy[types.GroupInfo](in)

}

func ToGroupInfoPbTypes(in *types.GroupInfo) *dm.GroupInfo {
	return utils.Copy[dm.GroupInfo](in)
}
