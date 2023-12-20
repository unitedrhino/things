package logic

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
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

func ToVidPageRpc(in *types.PageInfo) *vid.PageInfo {
	if in == nil {
		return nil
	}
	return &vid.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToDmPageRpc(in *types.PageInfo) *dm.PageInfo {
	if in == nil {
		return nil
	}
	return &dm.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToSysPageRpc(in *types.PageInfo) *sys.PageInfo {
	if in == nil {
		return nil
	}
	return &sys.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToTimedJobPageRpc(in *types.PageInfo) *timedjob.PageInfo {
	if in == nil {
		return nil
	}
	return &timedjob.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
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

func ToDiPageRpc(in *types.PageInfo) *dm.PageInfo {
	if in == nil {
		return nil
	}
	return &dm.PageInfo{
		Page: in.Page,
		Size: in.Size,
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

func ToDmPointApi(in *dm.Point) *types.Point {
	if in == nil {
		return nil
	}
	return &types.Point{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func ToSysPointRpc(in *types.Point) *sys.Point {
	if in == nil {
		return nil
	}
	return &sys.Point{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func ToSysPointApi(in *sys.Point) *types.Point {
	if in == nil {
		return nil
	}
	return &types.Point{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func ToOtaPageRpc(in *types.PageInfo) *dm.OtaPageInfo {
	if in == nil {
		return nil
	}
	return &dm.OtaPageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToDiSendOption(in *types.SendOption) *dm.SendOption {
	if in == nil {
		return nil
	}
	return &dm.SendOption{
		TimeoutToFail:  in.TimeoutToFail,
		RequestTimeout: in.RequestTimeout,
		RetryInterval:  in.RetryInterval,
	}
}

func SysToWithIDTypes(in *sys.WithID) *types.WithID {
	if in == nil {
		return nil
	}
	return &types.WithID{
		ID: in.Id,
	}
}
