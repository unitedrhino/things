package rulelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func ToSceneInfoDo(in *ud.SceneInfo) *scene.Info {
	if in == nil {
		return nil
	}
	do := utils.Copy[scene.Info](in)
	if do.Type == scene.SceneTypeAuto {
		do.If = utils.UnmarshalNoErr[scene.If](in.If)
		do.When = utils.UnmarshalNoErr[scene.When](in.When)
	}
	do.Then = utils.UnmarshalNoErr[scene.Then](in.Then)
	return do
}

func ToSceneInfoPo(in *scene.Info) *relationDB.UdSceneInfo {
	return &relationDB.UdSceneInfo{
		ID:        in.ID,
		Type:      in.Type,
		ProjectID: stores.ProjectID(in.ProjectID),
		AreaID:    stores.AreaID(in.AreaID),
		//AreaIDs: in.AreaIDs,
		FlowPath:    in.FlowPath,
		Name:        in.Name,
		Desc:        in.Desc,
		Tag:         in.Tag,
		Logo:        in.Logo,
		Body:        in.Body,
		DeviceMode:  in.DeviceMode,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		DeviceAlias: in.DeviceAlias,
		HeadImg:     in.HeadImg,
		Status:      in.Status,
		IsCommon:    in.IsCommon,
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
	return &relationDB.UdSceneIfTrigger{
		SceneID: si.ID,
		Type:    in.Type,
		Status:  si.Status,
		LastRunTime: sql.NullTime{
			Time:  in.Timer.GenLastRunTime(now),
			Valid: true,
		},
		Order:     in.Order,
		ProjectID: stores.ProjectID(si.ProjectID),
		AreaID:    stores.AreaID(si.AreaID),
		Device:    ToSceneTriggerDevicePo(in.Device),
		Timer:     ToSceneTriggerTimerPo(si, in.Timer),
		Weather:   utils.Copy2[relationDB.UdSceneTriggerWeather](in.Weather),
	}
}

func ToSceneTriggerTimerPo(si *scene.Info, in *scene.TriggerTimer) (ret relationDB.UdSceneTriggerTimer) {
	if in == nil {
		return relationDB.UdSceneTriggerTimer{}
	}
	return relationDB.UdSceneTriggerTimer{
		ExecAt:        in.ExecAt,
		ExecAdd:       in.ExecAdd,
		ExecRepeat:    utils.BStrToInt64(in.ExecRepeat),
		ExecType:      in.ExecType,
		ExecLoopStart: in.ExecLoopStart,
		ExecLoopEnd:   in.ExecLoopEnd,
		ExecLoop:      in.ExecLoop,
		RepeatType:    in.RepeatType,
	}
}

func ToSceneTriggerDevicePo(in *scene.TriggerDevice) (ret relationDB.UdSceneTriggerDevice) {
	if in == nil {
		return relationDB.UdSceneTriggerDevice{}
	}
	return relationDB.UdSceneTriggerDevice{
		ProductID:        in.ProductID,
		ProductName:      in.ProductName,
		SelectType:       in.SelectType,
		DeviceName:       in.DeviceName,
		DeviceAlias:      in.DeviceAlias,
		Type:             in.Type,
		DataID:           in.DataID,
		DataName:         in.DataName,
		TermType:         string(in.TermType),
		Terms:            in.Terms,
		Values:           in.Values,
		SchemaAffordance: in.SchemaAffordance,
		Body:             in.Body,
		StateKeep:        utils.Copy2[relationDB.UdStateKeep](in.StateKeep),
	}
}

func PoToSceneInfoDo(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.UdSceneInfo) *scene.Info {
	if in == nil {
		return nil
	}
	ret := &scene.Info{
		ID:          in.ID,
		ProjectID:   int64(in.ProjectID),
		AreaID:      int64(in.AreaID),
		Name:        in.Name,
		Tag:         in.Tag,
		Logo:        in.Logo,
		HeadImg:     in.HeadImg,
		FlowPath:    in.FlowPath,
		Desc:        in.Desc,
		Body:        in.Body,
		CreatedTime: in.CreatedTime,
		DeviceMode:  in.DeviceMode,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		DeviceAlias: in.DeviceAlias,
		LastRunTime: utils.GetNullTime(in.LastRunTime),
		Type:        in.Type,
		If: scene.If{
			Triggers: ToSceneTriggersDo(ctx, svcCtx, in.Triggers),
		},
		When: scene.When{
			ValidRanges:   in.UdSceneWhen.ValidRanges,
			InvalidRanges: in.UdSceneWhen.InvalidRanges,
			Conditions:    in.UdSceneWhen.Conditions,
		},
		Then: scene.Then{
			Actions: ToSceneActionsDo(ctx, svcCtx, in.UdSceneThen.Actions),
		},
		Status:   in.Status,
		IsCommon: in.IsCommon,
	}
	for i, v := range ret.When.Conditions.Terms {
		if v.ColumnType != scene.TermColumnTypeProperty {
			continue
		}
		p := v.Property
		di, err := svcCtx.DeviceCache.GetData(ctx, devices.Core{
			ProductID:  p.ProductID,
			DeviceName: p.DeviceName,
		})
		if err == nil {
			v.Property.ProductName = di.ProductName
			v.Property.DeviceAlias = di.DeviceAlias.GetValue()
			ret.When.Conditions.Terms[i] = v
		}
	}
	return ret
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
		//ID:      in.ID,
		Order:   in.Order,
		SceneID: s.ID,
		Type:    in.Type,
		Delay:   in.Delay,
		Notify:  in.Notify,
		Alarm:   in.Alarm,
	}
	if in.Scene != nil {
		po.Scene = relationDB.UdSceneActionScene{SceneID: in.Scene.SceneID, AreaID: in.Scene.AreaID, SceneType: in.Scene.SceneType, SceneName: in.Scene.SceneName}
	}
	if in.Device != nil {
		po.Device = relationDB.UdSceneActionDevice{
			//ProjectID:        int64(in.Device.ProjectID),
			AreaID:           in.Device.AreaID,
			AreaName:         in.Device.AreaName,
			ProductID:        in.Device.ProductID,
			ProductName:      in.Device.ProductName,
			SelectType:       in.Device.SelectType,
			DeviceName:       in.Device.DeviceName,
			DeviceAlias:      in.Device.DeviceAlias,
			DataName:         in.Device.DataName,
			GroupID:          in.Device.GroupID,
			Type:             in.Device.Type,
			DataID:           in.Device.DataID,
			Value:            in.Device.Value,
			SchemaAffordance: in.Device.SchemaAffordance,
			Values:           in.Device.Values,
			Body:             in.Device.Body,
		}
	}
	return po
}

func ToSceneActionsDo(ctx context.Context, svcCtx *svc.ServiceContext, in []*relationDB.UdSceneThenAction) (ret scene.Actions) {
	for _, v := range in {
		ret = append(ret, ToSceneActionDo(ctx, svcCtx, v))
	}
	return
}

func ToSceneActionDo(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.UdSceneThenAction) *scene.Action {
	if in == nil {
		return nil
	}
	do := &scene.Action{
		ID:     in.ID,
		Order:  in.Order,
		Type:   in.Type,
		Delay:  in.Delay,
		Notify: in.Notify,
		Alarm:  in.Alarm,
	}
	do.Scene = &scene.ActionScene{SceneID: in.Scene.SceneID, AreaID: in.Scene.AreaID, SceneType: in.Scene.SceneType, SceneName: in.Scene.SceneName}
	do.Device = &scene.ActionDevice{
		//ProjectID:        int64(in.Device.ProjectID),
		AreaID:           in.Device.AreaID,
		AreaName:         in.Device.AreaName,
		ProductID:        in.Device.ProductID,
		ProductName:      in.Device.ProductName,
		SelectType:       in.Device.SelectType,
		DeviceName:       in.Device.DeviceName,
		DeviceAlias:      in.Device.DeviceAlias,
		GroupID:          in.Device.GroupID,
		Type:             in.Device.Type,
		DataID:           in.Device.DataID,
		DataName:         in.Device.DataName,
		Value:            in.Device.Value,
		SchemaAffordance: in.Device.SchemaAffordance,
		Values:           in.Device.Values,
		Body:             in.Device.Body,
	}
	if in.Type == scene.ActionExecutorDevice {
		di, err := svcCtx.DeviceCache.GetData(ctx, devices.Core{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		})
		if err == nil {
			do.Device.ProductName = di.ProductName
			do.Device.DeviceAlias = di.DeviceAlias.GetValue()
		}
	}
	return do
}

func ToSceneTriggersDo(ctx context.Context, svcCtx *svc.ServiceContext, in []*relationDB.UdSceneIfTrigger) (ret scene.Triggers) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, ToSceneTriggerDo(ctx, svcCtx, v))
	}
	return
}

func ToSceneTriggerDo(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.UdSceneIfTrigger) *scene.Trigger {
	if in == nil {
		return nil
	}
	return &scene.Trigger{
		Type:    in.Type,
		Order:   in.Order,
		AreaID:  int64(in.AreaID),
		Device:  ToSceneTriggerDeviceDo(ctx, svcCtx, in.Device),
		Timer:   ToSceneTriggerTimerDo(in.Timer),
		Weather: utils.Copy[scene.TriggerWeather](in.Weather),
	}
}

func ToSceneTriggerTimerDo(in relationDB.UdSceneTriggerTimer) (ret *scene.TriggerTimer) {
	return &scene.TriggerTimer{
		ExecAt:        in.ExecAt,
		ExecAdd:       in.ExecAdd,
		ExecRepeat:    utils.Int64ToBStr(in.ExecRepeat, scene.RepeatTypeLen[in.RepeatType]),
		ExecType:      in.ExecType,
		ExecLoopStart: in.ExecLoopStart,
		ExecLoopEnd:   in.ExecLoopEnd,
		ExecLoop:      in.ExecLoop,
		RepeatType:    in.RepeatType,
	}
}

func ToSceneTriggerDeviceDo(ctx context.Context, svcCtx *svc.ServiceContext, in relationDB.UdSceneTriggerDevice) (ret *scene.TriggerDevice) {
	ret = &scene.TriggerDevice{
		ProductID:   in.ProductID,
		ProductName: in.ProductName,
		SelectType:  in.SelectType,
		DeviceName:  in.DeviceName,
		DeviceAlias: in.DeviceAlias,
		Type:        in.Type,
		Compare: scene.Compare{
			DataID:   in.DataID,
			DataName: in.DataName,
			TermType: scene.CmpType(in.TermType),
			Values:   in.Values,
			Terms:    in.Terms,
		},
		SchemaAffordance: in.SchemaAffordance,
		Body:             in.Body,
		StateKeep:        utils.Copy[scene.StateKeep](in.StateKeep),
	}
	di, err := svcCtx.DeviceCache.GetData(ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err == nil {
		ret.ProductName = di.ProductName
		ret.DeviceAlias = di.DeviceAlias.GetValue()
	}
	return
}

func PoToSceneInfoPb(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.UdSceneInfo) *ud.SceneInfo {
	if in == nil {
		return nil
	}
	do := PoToSceneInfoDo(ctx, svcCtx, in)
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
