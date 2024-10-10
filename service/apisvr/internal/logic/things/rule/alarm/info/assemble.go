package info

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"
)

func AlarmInfoToApi(in *ud.AlarmInfo) *types.AlarmInfo {
	return utils.Copy[types.AlarmInfo](in)
}
