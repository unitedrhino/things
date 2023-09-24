package auth

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToAreaApis(in []*sys.UserAuthArea) (ret []*types.UserAuthArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserAuthArea{AreaID: v.AreaID})
	}
	return
}

func ToAreaPbs(in []*types.UserAuthArea) (ret []*sys.UserAuthArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserAuthArea{AreaID: v.AreaID})
	}
	return
}

func ToProjectApis(in []*sys.UserAuthProject) (ret []*types.UserAuthProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserAuthProject{ProjectID: v.ProjectID})
	}
	return
}

func ToProjectPbs(in []*types.UserAuthProject) (ret []*sys.UserAuthProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserAuthProject{ProjectID: v.ProjectID})
	}
	return
}
