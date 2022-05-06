package logic

import (
	"context"

	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/i-Things/things/src/ddsvr/pb/dd"

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

// ping pong 测试
func (l *PingLogic) Ping(in *dd.Empty) (*dd.Empty, error) {
	// todo: add your logic here and delete this line

	return &dd.Empty{}, nil
}
