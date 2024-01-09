package rulelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/udsvr/internal/domain/scene"
	"github.com/i-Things/things/src/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/udsvr/pb/ud"
)

func ToSceneInfoDo(in *ud.SceneInfo) *scene.Info {
	if in == nil {
		return nil
	}

	return &scene.Info{
		ID:      in.Id,
		Name:    in.Name,
		Desc:    in.Desc,
		AreaID:  in.AreaID,
		Trigger: utils.UnmarshalNoErr[scene.Trigger](in.Trigger),
		When:    utils.UnmarshalNoErr[scene.When](in.When),
		Then:    utils.UnmarshalNoErr[scene.Then](in.Then),
		Status:  in.Status,
	}
}

func ToSceneInfoPo(in *scene.Info) *relationDB.UdSceneInfo {
	return &relationDB.UdSceneInfo{
		ID:     in.ID,
		AreaID: in.AreaID,
		Name:   in.Name,
		Desc:   in.Desc,
		UdSceneTrigger: relationDB.UdSceneTrigger{
			Type:   string(in.Trigger.Type),
			Device: in.Trigger.Device,
			Timer:  in.Trigger.Timer,
		},
		UdSceneWhen: relationDB.UdSceneWhen{
			TermCondType: string(in.When.TermCondType),
			ValidRange:   in.When.ValidRange,
			InvalidRange: in.When.InvalidRange,
			Terms:        in.When.Terms,
		},
		UdSceneThen: relationDB.UdSceneThen{Actions: in.Then.Actions},
	}
}

func PoToSceneInfoDo(in *relationDB.UdSceneInfo) *scene.Info {
	if in == nil {
		return nil
	}
	return &scene.Info{
		ID:          in.ID,
		AreaID:      in.AreaID,
		Name:        in.Name,
		Desc:        in.Desc,
		CreatedTime: in.CreatedTime,
		Trigger: scene.Trigger{
			Type:   scene.TriggerType(in.UdSceneTrigger.Type),
			Device: in.UdSceneTrigger.Device,
			Timer:  in.UdSceneTrigger.Timer,
		},
		When: scene.When{
			TermCondType: scene.TermCondType(in.UdSceneWhen.TermCondType),
			ValidRange:   in.UdSceneWhen.ValidRange,
			InvalidRange: in.UdSceneWhen.InvalidRange,
			Terms:        in.UdSceneWhen.Terms,
		},
		Then: scene.Then{
			Actions: in.UdSceneThen.Actions,
		},
		Status: in.Status,
	}
}

func PoToSceneInfoPb(in *relationDB.UdSceneInfo) *ud.SceneInfo {
	if in == nil {
		return nil
	}
	do := PoToSceneInfoDo(in)
	return &ud.SceneInfo{
		Id:      in.ID,
		Name:    in.Name,
		Desc:    in.Desc,
		AreaID:  in.AreaID,
		Trigger: utils.MarshalNoErr(do.Trigger),
		When:    utils.MarshalNoErr(do.When),
		Then:    utils.MarshalNoErr(do.Then),
		Status:  in.Status,
	}
}
func PoToSceneInfoPbs(in []*relationDB.UdSceneInfo) (ret []*ud.SceneInfo) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, PoToSceneInfoPb(v))
	}
	return ret
}
