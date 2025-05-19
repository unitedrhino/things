package deviceGroup

import (
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

const DictCode = "deviceGroupPurpose"

const DictDefault = "default"

var DefaultBuildingConfig = GroupConfig{
	UniqueDevice: true,
	TsIndex:      true,
}

type (
	GroupConfig struct {
		UniqueDevice bool `json:"uniqueDevice"` //一个设备只能存在该用途下的一个分组中
		TsIndex      bool `json:"tsIndex"`      //开启时序索引
	}
	GroupDetail struct {
		*GroupConfig `json:"config"`
		*sys.DictDetail
	}
)

func ToMapGroup(in []*sys.DictDetail) (ret map[string]*GroupDetail) {
	ret = make(map[string]*GroupDetail)
	for _, v := range in {
		if v.Body.GetValue() == "" {
			ret[v.Value] = &GroupDetail{DictDetail: v, GroupConfig: &GroupConfig{}}
			continue
		}
		var c GroupConfig
		err := json.Unmarshal([]byte(v.Body.GetValue()), &c)
		if err != nil {
			logx.Error(err)
			continue
		}
		ret[v.Value] = &GroupDetail{DictDetail: v, GroupConfig: &c}
	}
	return
}

func NewGroupDetail(in []*sys.DictDetail) (ret []*GroupDetail) {
	for _, v := range in {
		if v.Body.GetValue() == "" {
			continue
		}
		var c GroupConfig
		err := json.Unmarshal([]byte(v.Body.GetValue()), &c)
		if err != nil {
			logx.Error(err)
			continue
		}
		if !c.TsIndex {
			continue
		}
		ret = append(ret, &GroupDetail{DictDetail: v, GroupConfig: &c})
	}
	return
}
