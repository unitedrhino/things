package logic

import (
	"context"

	"yl/src/user/api/test/internal/svc"
	"yl/src/user/api/test/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type Register2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegister2Logic(ctx context.Context, svcCtx *svc.ServiceContext) Register2Logic {
	return Register2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Register2Logic) Register2(req types.Register2Req) error {
	// todo: add your logic here and delete this line

	return nil
}
