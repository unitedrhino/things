package timedschedulerlogic

import (
	"context"

	"github.com/i-Things/things/src/timed/timedschedulersvr/internal/svc"
	"github.com/i-Things/things/src/timed/timedschedulersvr/pb/timedscheduler"

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

func (l *PingLogic) Ping(in *timedscheduler.Empty) (*timedscheduler.Empty, error) {
	// todo: add your logic here and delete this line

	return &timedscheduler.Empty{}, nil
}
