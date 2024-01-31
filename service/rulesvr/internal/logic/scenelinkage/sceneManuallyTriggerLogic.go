package scenelinkagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/service/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/service/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneManuallyTriggerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	SiDB *relationDB.SceneInfoRepo
}

func NewSceneManuallyTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneManuallyTriggerLogic {
	return &SceneManuallyTriggerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		SiDB:   relationDB.NewSceneInfoRepo(ctx),
	}
}

func (l *SceneManuallyTriggerLogic) SceneManuallyTrigger(in *rule.WithID) (*rule.Empty, error) {
	pi, err := l.SiDB.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if pi.Status != def.True {
		return nil, errors.NotEnable
	}
	if pi.TriggerType != scene.TriggerTypeManual {
		return nil, errors.Parameter.AddMsg("只支持手动触发模式")
	}
	err = pi.Then.Execute(l.ctx, scene.ActionRepo{
		DeviceInteract: l.svcCtx.DeviceInteract,
		DeviceM:        l.svcCtx.DeviceM,
	})
	if err != nil {
		return nil, err
	}
	return &rule.Empty{}, nil
}
