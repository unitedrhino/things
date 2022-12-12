package logic

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/disvr/pb/di"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
	if in == nil {
		return nil
	}
	return &dm.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToDiPageRpc(in *types.PageInfo) *di.PageInfo {
	if in == nil {
		return nil
	}
	return &di.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}
