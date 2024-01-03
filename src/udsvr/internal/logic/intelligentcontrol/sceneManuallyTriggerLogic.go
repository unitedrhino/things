package intelligentcontrollogic

import (
	"context"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

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
	// todo: add your logic here and delete this line

	return &ud.Empty{}, nil
}
