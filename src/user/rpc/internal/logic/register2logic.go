package logic

import (
	"context"

	"yl/src/user/rpc/internal/svc"
	"yl/src/user/rpc/user"

	"github.com/tal-tech/go-zero/core/logx"
)

type Register2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Register2Logic {
	return &Register2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Register2Logic) Register2(in *user.Register2Req) (*user.Register2Resp, error) {
	// todo: add your logic here and delete this line

	return &user.Register2Resp{}, nil
}
