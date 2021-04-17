package logic

import (
	"context"

	"yl/src/webapi/api/test/internal/svc"
	"yl/src/webapi/api/test/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type RegisterCoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterCoreLogic {
	return RegisterCoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterCoreLogic) RegisterCore(req types.RegisterCoreReq) (*types.RegisterCoreResp, error) {
	// todo: add your logic here and delete this line

	return &types.RegisterCoreResp{}, nil
}
