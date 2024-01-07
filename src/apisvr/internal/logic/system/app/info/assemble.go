package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToAppInfoRpc(in *types.AppInfo) *sys.AppInfo {
	if in == nil {
		return nil
	}
	return &sys.AppInfo{
		Id:      in.ID,
		Name:    in.Name,
		Type:    in.Type,
		Code:    in.Code,
		Desc:    utils.ToRpcNullString(in.Desc),
		BaseUrl: in.BaseUrl,
		LogoUrl: in.LogoUrl,
	}
}

func ToAppInfoTypes(in *sys.AppInfo) *types.AppInfo {
	if in == nil {
		return nil
	}
	return &types.AppInfo{
		ID:      in.Id,
		Name:    in.Name,
		Type:    in.Type,
		Code:    in.Code,
		Desc:    utils.ToNullString(in.Desc),
		BaseUrl: in.BaseUrl,
		LogoUrl: in.LogoUrl,
	}
}

func ToAppInfosTypes(in []*sys.AppInfo) []*types.AppInfo {
	var ret []*types.AppInfo
	for _, v := range in {
		ret = append(ret, ToAppInfoTypes(v))
	}
	return ret
}
