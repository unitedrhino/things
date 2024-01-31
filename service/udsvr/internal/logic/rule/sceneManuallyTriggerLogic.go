package rulelogic

import (
	"context"
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
	if si.UdSceneTrigger.Type != string(scene.TriggerTypeManual) {
		return nil, errors.SceneTriggerType.AddMsg("该场景不是手动触发类型,无法执行")
	}
	//todo 执行触发
	return &ud.Empty{}, nil
}
