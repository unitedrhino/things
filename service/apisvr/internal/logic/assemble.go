package logic

import (
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
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

func ToPageResp(p *types.PageInfo, total int64) types.PageResp {
	ret := types.PageResp{Total: total}
	if p == nil {
		return ret
	}
	ret.Page = p.Page
	ret.PageSize = p.Size
	return ret
}
