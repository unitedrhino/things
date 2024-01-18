package area

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToAreaApis(in []*sys.UserArea) (ret []*types.UserArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserArea{AreaID: v.AreaID})
	}
	return
}

func ToAreaPbs(in []*types.UserArea) (ret []*sys.UserArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserArea{AreaID: v.AreaID})
	}
	return
}
