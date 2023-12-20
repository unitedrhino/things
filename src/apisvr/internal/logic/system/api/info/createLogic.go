package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
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

func (l *CreateLogic) Create(req *types.ApiInfo) (*types.WithID, error) {
	resp, err := l.svcCtx.ApiRpc.ApiInfoCreate(l.ctx, ToApiInfoRpc(req))
	if err != nil {
		err := errors.Fmt(err)
		l.Errorf("%s.rpc.ApiCreate req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	return logic.SysToWithIDTypes(resp), nil
}
