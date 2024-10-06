package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

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

func (l *CreateLogic) Create(req *types.SceneInfoCreateReq) (*types.CommonResp, error) {
	rst, err := l.svcCtx.Rule.SceneInfoCreate(l.ctx, ToScenePb(&req.SceneInfo))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.SceneInfoCreate req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.CommonResp{ID: rst.Id}, nil
}
