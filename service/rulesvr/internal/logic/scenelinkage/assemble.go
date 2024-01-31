package scenelinkagelogic

import (
	"encoding/json"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/service/rulesvr/pb/rule"
	"time"
)

func ToSceneDo(in *rule.SceneInfo) (*scene.Info, error) {
	var err error
	info := scene.Info{
		ID:          in.Id,
		Name:        in.Name,
		Desc:        in.Desc,
		Status:      in.Status,
		TriggerType: scene.TriggerType(in.TriggerType),
		CreatedTime: time.Unix(in.CreatedTime, 0),
	}
	switch info.TriggerType {
	case scene.TriggerTypeDevice, scene.TriggerTypeTimer:
		err = json.Unmarshal([]byte(in.Trigger), &info.Trigger)
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(in.When), &info.When)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(in.Then), &info.Then)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func ToScenePb(in *scene.Info) *rule.SceneInfo {
	ret := rule.SceneInfo{
		Id:          in.ID,
		Name:        in.Name,
		Desc:        in.Desc,
		Status:      in.Status,
		TriggerType: string(in.TriggerType),
		Trigger:     utils.MarshalNoErr(in.Trigger),
		When:        utils.MarshalNoErr(in.When),
		Then:        utils.MarshalNoErr(in.Then),
		CreatedTime: in.CreatedTime.Unix(),
	}
	switch in.TriggerType {
	case scene.TriggerTypeDevice:
		ret.Trigger = utils.MarshalNoErr(in.Trigger.Device)
	}

	return &ret
}
