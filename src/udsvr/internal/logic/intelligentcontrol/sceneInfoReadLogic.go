package intelligentcontrollogic

import (
	"context"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoReadLogic {
	return &SceneInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoReadLogic) SceneInfoRead(in *ud.WithID) (*ud.SceneInfo, error) {
	// todo: add your logic here and delete this line

	return &ud.SceneInfo{}, nil
}
