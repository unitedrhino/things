package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.SceneInfoCreateReq) error {
	_, err := l.svcCtx.Scene.SceneInfoCreate(l.ctx, ToScenePb(&req.SceneInfo))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneInfoCreate req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
