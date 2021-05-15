package logic

import (
	"context"

	"yl/src/dmsvr/dm"
	"yl/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *dm.Request) (*dm.Response, error) {
	// todo: add your logic here and delete this line

	return &dm.Response{}, nil
}
