package timing

import (
	"context"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TimingHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewServerHandle(ctx context.Context, svcCtx *svc.ServiceContext) *TimingHandle {
	return &TimingHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *TimingHandle) DeviceTiming() error {
	//l.Infof("DeviceTiming")
	return nil
}

func (l *TimingHandle) SceneTiming() error {
	//l.Infof("SceneTiming")
	return nil
}
