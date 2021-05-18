package logic

import (
	"context"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginAuthLogic {
	return &LoginAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginAuthLogic) LoginAuth(in *dm.LoginAuthReq) (*dm.Response, error) {
	// todo: add your logic here and delete this line

	return &dm.Response{}, nil
}
