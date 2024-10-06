package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/udsvr/internal/domain/scene"
	"gitee.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"time"

	"gitee.com/i-Things/things/service/udsvr/internal/svc"
	"gitee.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneManuallyTriggerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneManuallyTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneManuallyTriggerLogic {
	return &SceneManuallyTriggerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneManuallyTriggerLogic) SceneManuallyTrigger(in *ud.WithID) (*ud.Empty, error) {
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
	si, err := relationDB.NewSceneInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	//if si.Type != string(scene.SceneTypeManual) {
	//	return nil, errors.TriggerType.AddMsg("该场景不是手动触发类型,无法执行")
	//}
	do := PoToSceneInfoDo(si)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) { //异步执行
		err = stores.WithNoDebug(ctx, relationDB.NewSceneInfoRepo).UpdateWithField(ctx, relationDB.SceneInfoFilter{IDs: []int64{si.ID}}, map[string]any{"last_run_time": time.Now()})
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		err = do.Then.Execute(ctx, scene.ActionRepo{
			Info:           do,
			DeviceInteract: l.svcCtx.DeviceInteract,
			DeviceM:        l.svcCtx.DeviceM,
			ProductCache:   l.svcCtx.ProductCache,
			DeviceCache:    l.svcCtx.DeviceCache,
			DeviceG:        l.svcCtx.DeviceG,
			SceneExec: func(ctx context.Context, sceneID int64, status def.Bool) error {
				l.Error("not support yet")
				return nil
			},
			AlarmExec: func(ctx context.Context, in scene.AlarmSerial) error {
				l.Error("not support yet")
				return nil
			},
			SaveLog: func(ctx context.Context, log *scene.Log) error {
				po := utils.Copy[relationDB.UdSceneLog](log)
				err := stores.WithNoDebug(l.ctx, relationDB.NewSceneLogRepo).Insert(ctx, po)
				return err
			},
		})
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	})
	return &ud.Empty{}, nil
}
