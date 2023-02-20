package scenelinkagelogic

import (
	"encoding/json"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
)

func ToSceneDo(in *rule.SceneInfo) (*scene.Info, error) {
	var err error
	info := scene.Info{
		ID:   in.Id,
		Name: in.Name,
		Desc: in.Desc,
	}
	err = json.Unmarshal([]byte(in.Trigger), &info.Trigger)
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
	return &rule.SceneInfo{
		Id:      in.ID,
		Name:    in.Name,
		Desc:    in.Desc,
		Trigger: utils.MarshalNoErr(in.Trigger),
		When:    utils.MarshalNoErr(in.When),
		Then:    utils.MarshalNoErr(in.Then),
	}
}
