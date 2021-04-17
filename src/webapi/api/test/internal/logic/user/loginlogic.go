package logic

import (
	"context"

	"yl/src/webapi/api/test/internal/svc"
	"yl/src/webapi/api/test/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &types.LoginResp{}, nil
}
