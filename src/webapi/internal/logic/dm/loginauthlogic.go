package logic

import (
	"context"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginAuthLogic {
	return LoginAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginAuthLogic) LoginAuth(req types.LoginAuthReq) error {
	// todo: add your logic here and delete this line

	return nil
}
