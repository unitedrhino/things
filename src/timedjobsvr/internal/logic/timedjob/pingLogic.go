package timedjoblogic

import (
	"context"

	"github.com/i-Things/things/src/timedjobsvr/internal/svc"
	"github.com/i-Things/things/src/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *PingLogic) Ping(in *timedjob.Request) (*timedjob.Response, error) {
	// todo: add your logic here and delete this line

	return &timedjob.Response{}, nil
}
