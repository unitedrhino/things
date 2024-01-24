package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToModuleInfoRpc(in *types.ModuleInfo) *sys.ModuleInfo {
	return &sys.ModuleInfo{
		Id:         in.ID,
		Code:       in.Code,
		Name:       in.Name,
		Type:       in.Type,
		SubType:    in.SubType,
		Path:       in.Path,
		Desc:       utils.ToRpcNullString(in.Desc),
		Icon:       in.Icon,
		Url:        in.Url,
		Order:      in.Order,
		HideInMenu: in.HideInMenu,
		Body:       utils.ToRpcNullString(in.Body),
	}
}
func ToModuleInfoApi(in *sys.ModuleInfo) *types.ModuleInfo {
	if in == nil {
		return nil
	}
	return &types.ModuleInfo{
		ID:         in.Id,
		Code:       in.Code,
		Name:       in.Name,
		Type:       in.Type,
		SubType:    in.SubType,
		Path:       in.Path,
		Desc:       utils.ToNullString(in.Desc),
		Icon:       in.Icon,
		Url:        in.Url,
		Order:      in.Order,
		HideInMenu: in.HideInMenu,
		Body:       utils.ToNullString(in.Body),
	}

}
func ToModuleInfosApi(in []*sys.ModuleInfo) (ret []*types.ModuleInfo) {
	for _, v := range in {
		v1 := ToModuleInfoApi(v)
		if v1 != nil {
			ret = append(ret, v1)
		}
	}
	return
}
