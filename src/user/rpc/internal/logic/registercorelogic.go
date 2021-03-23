package logic

import (
	"context"

	"yl/src/user/rpc/internal/svc"
	"yl/src/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type RegisterCoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterCoreLogic {
	return &RegisterCoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterCoreLogic) RegisterCore(in *user.RegisterCoreReq) (*user.RegisterCoreResp, error) {
	// todo: add your logic here and delete this line

	return &user.RegisterCoreResp{}, nil
}
