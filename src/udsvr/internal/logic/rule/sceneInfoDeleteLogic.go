package rulelogic

import (
	"context"
	"github.com/i-Things/things/src/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type SceneInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSceneInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SceneInfoDeleteLogic {
	return &SceneInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SceneInfoDeleteLogic) SceneInfoDelete(in *ud.WithID) (*ud.Empty, error) {
	err := relationDB.NewSceneInfoRepo(l.ctx).Delete(l.ctx, in.Id)
	return &ud.Empty{}, err
}
