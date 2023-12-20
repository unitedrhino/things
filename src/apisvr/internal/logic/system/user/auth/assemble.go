package auth

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToAreaApis(in []*sys.UserArea) (ret []*types.UserAuthArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserAuthArea{AreaID: v.AreaID})
	}
	return
}

func ToAreaPbs(in []*types.UserAuthArea) (ret []*sys.UserArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserArea{AreaID: v.AreaID})
	}
	return
}

func ToProjectApis(in []*sys.UserProject) (ret []*types.UserAuthProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserAuthProject{ProjectID: v.ProjectID})
	}
	return
}

func ToProjectPbs(in []*types.UserAuthProject) (ret []*sys.UserProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserProject{ProjectID: v.ProjectID})
	}
	return
}
