package sceneChangeEvent

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type Handle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewHandle(ctx context.Context, svcCtx *svc.ServiceContext) *Handle {
	return &Handle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithRoot(ctx),
	}
}
func (l *Handle) SceneAreaDelete(areaIDs []int64) error {
	var ids []int64
	for _, areaID := range areaIDs {
		if areaID > def.NotClassified {
			ids = append(ids, areaID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	{
		sis, err := relationDB.NewSceneInfoRepo(l.ctx).FindByFilter(l.ctx,
			relationDB.SceneInfoFilter{AreaIDs: ids}, nil)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
			return err
		}
		for _, si := range sis {
			err := relationDB.NewSceneInfoRepo(l.ctx).Delete(l.ctx, si.ID)
			if err != nil {
				logx.WithContext(l.ctx).Error(err)
			}
		}
	}
	var sceneIDs []int64
	var reason = scene.ReasonAreaDelete
	{ //将包含该设备的场景状态修改为异常
		f := relationDB.SceneIfTriggerFilter{
			Statuses: []scene.Status{scene.StatusNormal, scene.StatusForbidden},
			AreaID:   stores.CmpIn(ids),
		}
		pos, err := relationDB.NewSceneIfTriggerRepo(l.ctx).FindByFilter(l.ctx, f, nil)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		} else {
			for _, si := range pos {
				sceneIDs = append(sceneIDs, si.SceneID)
			}
		}
		err = relationDB.NewSceneIfTriggerRepo(l.ctx).UpdateWithField(l.ctx, f,
			map[string]any{"status": scene.StatusAbnormal, "reason": reason})
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	{
		f := relationDB.SceneActionFilter{
			DeviceAreaID:     stores.CmpIn(ids),
			DeviceSelectType: scene.SelectArea,
			Statuses:         []scene.Status{scene.StatusNormal, scene.StatusForbidden},
		}
		pos, err := relationDB.NewSceneActionRepo(l.ctx).FindByFilter(l.ctx, f, nil)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		} else {
			for _, si := range pos {
				sceneIDs = append(sceneIDs, si.SceneID)
			}
		}
		err = relationDB.NewSceneActionRepo(l.ctx).UpdateWithField(l.ctx, f,
			map[string]any{"status": scene.StatusAbnormal, "reason": scene.ReasonDeviceDelete})
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	if len(sceneIDs) > 0 {
		err := relationDB.NewSceneInfoRepo(l.ctx).UpdateWithField(l.ctx, relationDB.SceneInfoFilter{IDs: sceneIDs},
			map[string]any{"status": scene.StatusAbnormal, "reason": scene.ReasonDeviceDelete})
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	return nil
}

func (l *Handle) SceneProjectDelete(pi int64) error {
	if pi == 0 {
		return nil
	}
	sis, err := relationDB.NewSceneInfoRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.SceneInfoFilter{ProjectID: pi}, nil)
	if err != nil {
		logx.WithContext(l.ctx).Error(err)
		return err
	}
	for _, si := range sis {
		err := relationDB.NewSceneInfoRepo(l.ctx).Delete(l.ctx, si.ID)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	return nil
}

func (l *Handle) SceneDeviceDelete(di devices.Core) error {
	{ //删除单设备定时
		if di.ProductID == "" || di.DeviceName == "" {
			return nil
		}
		sis, err := relationDB.NewSceneInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.SceneInfoFilter{ProductID: di.ProductID, DeviceName: di.DeviceName}, nil)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
			return err
		}
		for _, si := range sis {
			err := relationDB.NewSceneInfoRepo(l.ctx).Delete(l.ctx, si.ID)
			if err != nil {
				logx.WithContext(l.ctx).Error(err)
			}
		}
	}
	var sceneIDs []int64
	var reason = scene.ReasonDeviceDelete
	{ //将包含该设备的场景状态修改为异常
		f := relationDB.SceneIfTriggerFilter{
			Statuses: []scene.Status{scene.StatusNormal, scene.StatusForbidden},
			Device:   &devices.Core{ProductID: di.ProductID, DeviceName: di.DeviceName},
		}
		pos, err := relationDB.NewSceneIfTriggerRepo(l.ctx).FindByFilter(l.ctx, f, nil)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		} else {
			for _, si := range pos {
				sceneIDs = append(sceneIDs, si.SceneID)
			}
		}
		err = relationDB.NewSceneIfTriggerRepo(l.ctx).UpdateWithField(l.ctx, f,
			map[string]any{"status": scene.StatusAbnormal, "reason": reason})
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	{
		f := relationDB.SceneActionFilter{
			ProductID:        di.ProductID,
			DeviceName:       di.DeviceName,
			DeviceSelectType: scene.SelectDeviceFixed,
			Statuses:         []scene.Status{scene.StatusNormal, scene.StatusForbidden},
		}
		pos, err := relationDB.NewSceneActionRepo(l.ctx).FindByFilter(l.ctx, f, nil)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		} else {
			for _, si := range pos {
				sceneIDs = append(sceneIDs, si.SceneID)
			}
		}
		err = relationDB.NewSceneActionRepo(l.ctx).UpdateWithField(l.ctx, f,
			map[string]any{"status": scene.StatusAbnormal, "reason": scene.ReasonDeviceDelete})
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	if len(sceneIDs) > 0 {
		err := relationDB.NewSceneInfoRepo(l.ctx).UpdateWithField(l.ctx, relationDB.SceneInfoFilter{IDs: sceneIDs},
			map[string]any{"status": scene.StatusAbnormal, "reason": scene.ReasonDeviceDelete})
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
		}
	}
	return nil
}
