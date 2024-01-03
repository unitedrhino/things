package intelligentcontrollogic

import (
	"context"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoCreateLogic {
	return &SceneInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 场景
func (l *SceneInfoCreateLogic) SceneInfoCreate(in *ud.SceneInfo) (*ud.WithID, error) {
	// todo: add your logic here and delete this line

	return &ud.WithID{}, nil
}
