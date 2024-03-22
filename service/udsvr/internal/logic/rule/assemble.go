package rulelogic

import (
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/domain"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/udsvr/pb/ud"
	"time"
)

func ToSceneInfoDo(in *ud.SceneInfo) *scene.Info {
	if in == nil {
		return nil
	}

	return &scene.Info{
		ID:      in.Id,
		Name:    in.Name,
		HeadImg: in.HeadImg,
		Tag:     in.Tag,
		Desc:    in.Desc,
		Trigger: utils.UnmarshalNoErr[scene.Trigger](in.Trigger),
		When:    utils.UnmarshalNoErr[scene.When](in.When),
		Then:    utils.UnmarshalNoErr[scene.Then](in.Then),
		Status:  in.Status,
	}
}

func ToSceneInfoPo(in *scene.Info) *relationDB.UdSceneInfo {
	return &relationDB.UdSceneInfo{
		ID: in.ID,
		//AreaIDs: in.AreaIDs,
		Name:    in.Name,
		Desc:    in.Desc,
		Tag:     in.Tag,
		HeadImg: in.HeadImg,
		UdSceneTrigger: relationDB.UdSceneTrigger{
			Type: string(in.Trigger.Type),
			//Devices: ToSceneTriggerDevicesPo(in, in.Trigger.Devices),
			Timers: ToSceneTriggerTimersPo(in, in.Trigger.Timers),
		},
		UdSceneWhen: relationDB.UdSceneWhen{
			ValidRanges:   in.When.ValidRanges,
			InvalidRanges: in.When.InvalidRanges,
			Conditions:    in.When.Conditions,
		},
		UdSceneThen: relationDB.UdSceneThen{Actions: ToSceneActionsPo(in, in.Then.Actions)},
	}
}

func ToSceneTriggerTimersPo(si *scene.Info, in scene.Timers) (ret []*relationDB.UdSceneTriggerTimer) {
	if in == nil {
		return nil
	}
	now := time.Now()
	for _, v := range in {
		ret = append(ret, &relationDB.UdSceneTriggerTimer{
			SceneID:     si.ID,
			Status:      si.Status,
			ExecAt:      v.ExecAt,
			ExecRepeat:  v.ExecRepeat,
			LastRunTime: domain.GenLastRunTime(now, v.ExecAt),
		})
	}
	return
}

//func ToSceneTriggerDevicesPo(si *scene.Info, in scene.TriggerDevices) (ret []*relationDB.UdSceneTriggerDevice) {
//	if in == nil {
//		return nil
//	}
//	for _, v := range in {
//		ret = append(ret, &relationDB.UdSceneTriggerDevice{
//			SceneID:     si.ID,
//			ProductID:   v.ProductID,
//			SelectType:  string(v.SelectType),
//			DeviceName:  v.DeviceName,
//			DeviceAlias: v.DeviceAlias,
//			Type:        string(v.Type),
//			//Schema:      v.Schema,
//		})
//	}
//	return
//}

func PoToSceneInfoDo(in *relationDB.UdSceneInfo) *scene.Info {
	if in == nil {
		return nil
	}
	return &scene.Info{
		ID: in.ID,
		//AreaIDs:     in.AreaIDs,
		Name:        in.Name,
		Tag:         in.Tag,
		HeadImg:     in.HeadImg,
		Desc:        in.Desc,
		CreatedTime: in.CreatedTime,
		Trigger: scene.Trigger{
			Type:    scene.TriggerType(in.UdSceneTrigger.Type),
			Devices: ToSceneTriggerDevicesDo(in.UdSceneTrigger.Devices),
			Timers:  ToSceneTriggerTimersDo(in.UdSceneTrigger.Timers),
		},
		When: scene.When{
			ValidRanges:   in.UdSceneWhen.ValidRanges,
			InvalidRanges: in.UdSceneWhen.InvalidRanges,
			Conditions:    in.UdSceneWhen.Conditions,
		},
		Then: scene.Then{
			Actions: ToSceneActionsDo(in.UdSceneThen.Actions),
		},
		Status: in.Status,
	}
}

func ToSceneActionsPo(s *scene.Info, in scene.Actions) (ret []*relationDB.UdSceneThenAction) {
	for _, v := range in {
		ret = append(ret, ToSceneActionPo(s, v))
	}
	return
}

func ToSceneActionPo(s *scene.Info, in *scene.Action) *relationDB.UdSceneThenAction {
	if in == nil {
		return nil
	}
	po := &relationDB.UdSceneThenAction{
		ID:          in.ID,
		SceneID:     s.ID,
		ExecuteType: in.ExecuteType,
		Delay:       in.Delay,
	}
	if in.Device != nil {
		po.Device = &relationDB.UdSceneActionDevice{
			ProjectID:   stores.ProjectID(in.Device.ProjectID),
			AreaID:      stores.AreaID(in.Device.AreaID),
			ProductID:   in.Device.ProductID,
			SelectType:  in.Device.SelectType,
			DeviceName:  in.Device.DeviceName,
			DeviceAlias: in.Device.DeviceAlias,
			DataName:    in.Device.DataName,
			GroupID:     in.Device.GroupID,
			Type:        in.Device.Type,
			DataID:      in.Device.DataID,
			Value:       in.Device.Value,
		}
	}
	return po
}

func ToSceneActionsDo(in []*relationDB.UdSceneThenAction) (ret scene.Actions) {
	for _, v := range in {
		ret = append(ret, ToSceneActionDo(v))
	}
	return
}

func ToSceneActionDo(in *relationDB.UdSceneThenAction) *scene.Action {
	if in == nil {
		return nil
	}
	do := &scene.Action{
		ID:          in.ID,
		ExecuteType: in.ExecuteType,
		Delay:       in.Delay,
	}
	if in.Device != nil {
		do.Device = &scene.ActionDevice{
			ProjectID:   int64(in.Device.ProjectID),
			AreaID:      int64(in.Device.AreaID),
			ProductID:   in.Device.ProductID,
			SelectType:  in.Device.SelectType,
			DeviceName:  in.Device.DeviceName,
			DeviceAlias: in.Device.DeviceAlias,
			GroupID:     in.Device.GroupID,
			Type:        in.Device.Type,
			DataID:      in.Device.DataID,
			DataName:    in.Device.DataName,
			Value:       in.Device.Value,
		}
	}
	return do
}

func ToSceneTriggerTimersDo(in []*relationDB.UdSceneTriggerTimer) (ret scene.Timers) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, &scene.Timer{
			ExecAt:     v.ExecAt,
			ExecRepeat: v.ExecRepeat,
		})
	}
	return
}

func ToSceneTriggerDevicesDo(in []*relationDB.UdSceneTriggerDevice) (ret scene.TriggerDevices) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, &scene.TriggerDevice{
			ProductID:   v.ProductID,
			SelectType:  scene.SelectType(v.SelectType),
			DeviceName:  v.DeviceName,
			DeviceAlias: v.DeviceAlias,
			Type:        scene.TriggerDeviceType(v.Type),
			Schema:      ToSceneTriggerDeviceSchemaDo(v.Schema),
		})
	}
	return
}
func ToSceneTriggerDeviceSchemaDo(in *relationDB.UdSceneTriggerDeviceSchema) *scene.TriggerDeviceSchema {
	if in == nil {
		return nil
	}
	return &scene.TriggerDeviceSchema{
		DataID:   in.DataID,
		DataName: in.DataName,
		TermType: scene.CmpType(in.TermType),
		Values:   in.Values,
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
		Tag:     in.Tag,
		HeadImg: in.HeadImg,
		//AreaIDs: in.AreaIDs,
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

func ToDeviceTimerPb(in *relationDB.UdDeviceTimerInfo) *ud.DeviceTimerInfo {
	if in == nil {
		return nil
	}
	return &ud.DeviceTimerInfo{
		Id:   in.ID,
		Name: in.Name,
		Device: &ud.DeviceCore{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		},
		CreatedTime: in.CreatedTime.Unix(),
		TriggerType: in.TriggerType,
		ExecAt:      in.ExecAt,
		ExecRepeat:  in.ExecRepeat,
		ActionType:  in.ActionType,
		DataID:      in.DataID,
		Value:       in.Value,
		Status:      in.Status,
	}
}

func ToDeviceTimersPb(in []*relationDB.UdDeviceTimerInfo) (ret []*ud.DeviceTimerInfo) {
	for _, v := range in {
		ret = append(ret, ToDeviceTimerPb(v))
	}
	return
}
