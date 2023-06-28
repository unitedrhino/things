package info

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
)

func AlarmInfoToApi(in *rule.AlarmInfo) *types.AlarmInfo {
	pi := &types.AlarmInfo{
		ID:          in.Id,
		Name:        in.Name,
		Status:      in.Status,
		Desc:        in.Desc,
		CreatedTime: in.CreatedTime,
		Level:       in.Level,
	}
	return pi
}
