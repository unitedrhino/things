package custom

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToCustomTopicPb(info *types.ProductCustomTopic) *dm.CustomTopic {
	if info == nil {
		return nil
	}
	return &dm.CustomTopic{Topic: info.Topic, Direction: info.Direction}
}

func ToCustomTopicsPb(info []*types.ProductCustomTopic) (ret []*dm.CustomTopic) {
	if info == nil {
		return nil
	} else if len(info) == 0 {
		return []*dm.CustomTopic{}
	}
	for _, v := range info {
		ret = append(ret, ToCustomTopicPb(v))
	}
	return
}

func ToCustomTopicTypes(info *dm.CustomTopic) *types.ProductCustomTopic {
	if info == nil {
		return nil
	}
	return &types.ProductCustomTopic{Topic: info.Topic, Direction: info.Direction}
}

func ToCustomTopicsTypes(info []*dm.CustomTopic) (ret []*types.ProductCustomTopic) {
	for _, v := range info {
		ret = append(ret, ToCustomTopicTypes(v))
	}
	return
}
