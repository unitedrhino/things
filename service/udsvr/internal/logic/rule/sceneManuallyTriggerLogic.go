package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

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
	si, err := relationDB.NewSceneInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if si.UdSceneIf.Type != string(scene.IfTypeManual) {
		return nil, errors.TriggerType.AddMsg("该场景不是手动触发类型,无法执行")
	}
	do := PoToSceneInfoDo(si)
	ctxs.GoNewCtx(l.ctx, func(ctx context.Context) { //异步执行
		err = do.Then.Execute(ctx, scene.ActionRepo{
			DeviceInteract: l.svcCtx.DeviceInteract,
			DeviceM:        l.svcCtx.DeviceM,
			DeviceG:        l.svcCtx.DeviceG,
		})
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	})
	return &ud.Empty{}, nil
}
