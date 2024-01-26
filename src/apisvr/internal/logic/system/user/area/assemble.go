package area

import (
	"github.com/i-Things/things/src/apisvr/internal/logic/system/area/info"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToUserAreaDetail(in []*sys.UserArea, areaMap map[int64]*sys.AreaInfo) (ret []*types.UserAreaDetail) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.UserAreaDetail{AuthType: v.AuthType, AreaInfo: info.ToAreaInfoTypes(areaMap[v.AreaID], nil)})
	}
	return
}

func ToAreaPbs(in []*types.UserArea) (ret []*sys.UserArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.UserArea{AreaID: v.AreaID, AuthType: v.AuthType})
	}
	return
}
