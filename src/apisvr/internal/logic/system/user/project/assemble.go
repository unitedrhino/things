package project

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToProjectApis(in []*sys.UserProject) (ret []*types.UserProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserProject{ProjectID: v.ProjectID})
	}
	return
}

func ToProjectPbs(in []*types.UserProject) (ret []*sys.UserProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserProject{ProjectID: v.ProjectID})
	}
	return
}
