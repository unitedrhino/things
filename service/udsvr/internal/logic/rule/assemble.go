package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/domain"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func ToSceneInfoDo(in *ud.SceneInfo) *scene.Info {
	if in == nil {
		return nil
	}
	do := utils.Copy[scene.Info](in)
	do.If = utils.UnmarshalNoErr[scene.If](in.If)
	do.When = utils.UnmarshalNoErr[scene.When](in.When)
	do.Then = utils.UnmarshalNoErr[scene.Then](in.Then)
	return do
}

func ToSceneInfoPo(in *scene.Info) *relationDB.UdSceneInfo {
	return &relationDB.UdSceneInfo{
		ID:   in.ID,
		Type: in.Type,
		//AreaIDs: in.AreaIDs,
		FlowPath:    in.FlowPath,
		Name:        in.Name,
		Desc:        in.Desc,
		Tag:         in.Tag,
		DeviceMode:  in.DeviceMode,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		DeviceAlias: in.DeviceAlias,
		HeadImg:     in.HeadImg,
		Status:      in.Status,
		UdSceneIf: relationDB.UdSceneIf{
			Triggers: ToSceneTriggersPo(in, in.If.Triggers),
		},
		UdSceneWhen: relationDB.UdSceneWhen{
			ValidRanges:   in.When.ValidRanges,
			InvalidRanges: in.When.InvalidRanges,
			Conditions:    in.When.Conditions,
		},
		UdSceneThen: relationDB.UdSceneThen{Actions: ToSceneActionsPo(in, in.Then.Actions)},
	}
}

func ToSceneTriggersPo(si *scene.Info, in scene.Triggers) (ret []*relationDB.UdSceneIfTrigger) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, ToSceneTriggerPo(si, v))
	}
	return
}

func ToSceneTriggerPo(si *scene.Info, in *scene.Trigger) *relationDB.UdSceneIfTrigger {
	if in == nil {
		return nil
	}
	now := time.Now()
	var execAt int64
	if in.Timer != nil {
		execAt = in.Timer.ExecAt
	}
	return &relationDB.UdSceneIfTrigger{
		SceneID:     si.ID,
		Type:        in.Type,
		Status:      si.Status,
		LastRunTime: domain.GenLastRunTime(now, execAt),
		Order:       in.Order,
		Device:      ToSceneTriggerDevicePo(in.Device),
		Timer:       ToSceneTriggerTimerPo(si, in.Timer),
	}
}

func ToSceneTriggerTimerPo(si *scene.Info, in *scene.Timer) (ret relationDB.UdSceneTriggerTimer) {
	if in == nil {
		return relationDB.UdSceneTriggerTimer{}
	}
	return relationDB.UdSceneTriggerTimer{
		ExecAt:     in.ExecAt,
		ExecRepeat: in.ExecRepeat,
	}
}

func ToSceneTriggerDevicePo(in *scene.TriggerDevice) (ret relationDB.UdSceneTriggerDevice) {
	if in == nil {
		return relationDB.UdSceneTriggerDevice{}
	}
	return relationDB.UdSceneTriggerDevice{
		ProductID:        in.ProductID,
		SelectType:       in.SelectType,
		DeviceName:       in.DeviceName,
		DeviceAlias:      in.DeviceAlias,
		Type:             string(in.Type),
		DataID:           in.DataID,
		DataName:         in.DataName,
		TermType:         string(in.TermType),
		Values:           in.Values,
		SchemaAffordance: in.SchemaAffordance,
	}
}

func PoToSceneInfoDo(in *relationDB.UdSceneInfo) *scene.Info {
	if in == nil {
		return nil
	}
	return &scene.Info{
		ID:          in.ID,
		Name:        in.Name,
		Tag:         in.Tag,
		HeadImg:     in.HeadImg,
		FlowPath:    in.FlowPath,
		Desc:        in.Desc,
		CreatedTime: in.CreatedTime,
		DeviceMode:  in.DeviceMode,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		DeviceAlias: in.DeviceAlias,
		Type:        in.Type,
		If: scene.If{
			Triggers: ToSceneTriggersDo(in.Triggers),
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
		ID:      in.ID,
		SceneID: s.ID,
		Type:    in.Type,
		Delay:   in.Delay,
		Notify:  in.Notify,
	}
	if in.Device != nil {
		po.Device = relationDB.UdSceneActionDevice{
			//ProjectID:        int64(in.Device.ProjectID),
			//AreaID:           int64(in.Device.AreaID),
			ProductID:        in.Device.ProductID,
			SelectType:       in.Device.SelectType,
			DeviceName:       in.Device.DeviceName,
			DeviceAlias:      in.Device.DeviceAlias,
			DataName:         in.Device.DataName,
			GroupID:          in.Device.GroupID,
			Type:             in.Device.Type,
			DataID:           in.Device.DataID,
			Value:            in.Device.Value,
			SchemaAffordance: in.Device.SchemaAffordance,
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
		ID:    in.ID,
		Order: in.Order,
		Type:  in.Type,
		Delay: in.Delay,
	}
	do.Device = &scene.ActionDevice{
		//ProjectID:        int64(in.Device.ProjectID),
		//AreaID:           int64(in.Device.AreaID),
		ProductID:        in.Device.ProductID,
		SelectType:       in.Device.SelectType,
		DeviceName:       in.Device.DeviceName,
		DeviceAlias:      in.Device.DeviceAlias,
		GroupID:          in.Device.GroupID,
		Type:             in.Device.Type,
		DataID:           in.Device.DataID,
		DataName:         in.Device.DataName,
		Value:            in.Device.Value,
		SchemaAffordance: in.Device.SchemaAffordance,
	}

	return do
}

func ToSceneTriggersDo(in []*relationDB.UdSceneIfTrigger) (ret scene.Triggers) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, ToSceneTriggerDo(v))
	}
	return
}

func ToSceneTriggerDo(in *relationDB.UdSceneIfTrigger) *scene.Trigger {
	if in == nil {
		return nil
	}
	return &scene.Trigger{
		Type:   scene.TriggerType(in.Type),
		Order:  in.Order,
		Device: ToSceneTriggerDeviceDo(in.Device),
		Timer:  ToSceneTriggerTimerDo(in.Timer),
	}
}

func ToSceneTriggerTimerDo(in relationDB.UdSceneTriggerTimer) (ret *scene.Timer) {
	return &scene.Timer{
		ExecAt:     in.ExecAt,
		ExecRepeat: in.ExecRepeat,
	}
}

func ToSceneTriggerDeviceDo(in relationDB.UdSceneTriggerDevice) (ret *scene.TriggerDevice) {
	return &scene.TriggerDevice{
		ProductID:        in.ProductID,
		SelectType:       scene.SelectType(in.SelectType),
		DeviceName:       in.DeviceName,
		DeviceAlias:      in.DeviceAlias,
		Type:             scene.TriggerDeviceType(in.Type),
		DataID:           in.DataID,
		DataName:         in.DataName,
		TermType:         scene.CmpType(in.TermType),
		Values:           in.Values,
		SchemaAffordance: in.SchemaAffordance,
	}
}

func PoToSceneInfoPb(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.UdSceneInfo) *ud.SceneInfo {
	if in == nil {
		return nil
	}
	do := PoToSceneInfoDo(in)
	if in.HeadImg != "" {
		var err error
		in.HeadImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, in.HeadImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	pb := utils.Copy[ud.SceneInfo](in)
	pb.If = utils.MarshalNoErr(do.If)
	pb.When = utils.MarshalNoErr(do.When)
	pb.Then = utils.MarshalNoErr(do.Then)
	return pb
}

func PoToSceneInfoPbs(ctx context.Context, svcCtx *svc.ServiceContext, in []*relationDB.UdSceneInfo) (ret []*ud.SceneInfo) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, PoToSceneInfoPb(ctx, svcCtx, v))
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
