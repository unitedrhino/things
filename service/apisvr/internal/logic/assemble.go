package logic

import (
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/i-Things/things/service/rulesvr/pb/rule"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToTagsMap(tags []*types.Tag) map[string]string {
	if tags == nil {
		return nil
	}
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return tagMap
}

func ToTagsType(tags map[string]string) (retTag []*types.Tag) {
	for k, v := range tags {
		retTag = append(retTag, &types.Tag{
			Key:   k,
			Value: v,
		})
	}
	return
}

func ToDmPageRpc(in *types.PageInfo) *dm.PageInfo {
	return utils.Copy[dm.PageInfo](in)
}

func ToUdPageRpc(in *types.PageInfo) *ud.PageInfo {
	return utils.Copy[ud.PageInfo](in)
}

func ToRulePageRpc(in *types.PageInfo) *rule.PageInfo {
	if in == nil {
		return nil
	}
	return &rule.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}
func ToRuleTimeRangeRpc(in *types.TimeRange) *rule.TimeRange {
	if in == nil {
		return nil
	}
	return &rule.TimeRange{
		Start: in.Start,
		End:   in.End,
	}
}

func ToUdTimeRangeRpc(in *types.TimeRange) *ud.TimeRange {
	if in == nil {
		return nil
	}
	return &ud.TimeRange{
		Start: in.Start,
		End:   in.End,
	}
}

func ToDmPointRpc(in *types.Point) *dm.Point {
	if in == nil {
		return nil
	}
	return &dm.Point{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func ToDmSendOption(in *types.SendOption) *dm.SendOption {
	if in == nil {
		return nil
	}
	return &dm.SendOption{
		TimeoutToFail:  in.TimeoutToFail,
		RequestTimeout: in.RequestTimeout,
		RetryInterval:  in.RetryInterval,
	}
}
