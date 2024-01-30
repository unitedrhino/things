package area

import (
	"github.com/i-Things/things/src/apisvr/internal/logic/system/area/info"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToDataAreaDetail(in []*sys.DataArea, areaMap map[int64]*sys.AreaInfo) (ret []*types.DataAreaDetail) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.DataAreaDetail{AuthType: v.AuthType, AreaInfo: info.ToAreaInfoTypes(areaMap[v.AreaID], nil)})
	}
	return
}

func ToAreaPbs(in []*types.DataArea) (ret []*sys.DataArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.DataArea{AreaID: v.AreaID, AuthType: v.AuthType})
	}
	return
}
