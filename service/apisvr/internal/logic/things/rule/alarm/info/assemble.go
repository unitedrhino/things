package info

import (
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"gitee.com/i-Things/things/service/udsvr/pb/ud"
)

func AlarmInfoToApi(in *ud.AlarmInfo) *types.AlarmInfo {
	return utils.Copy[types.AlarmInfo](in)
}
