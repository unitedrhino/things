package staticEvent

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OneMinuteHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewOneMinuteHandle(ctx context.Context, svcCtx *svc.ServiceContext) *OneMinuteHandle {
	return &OneMinuteHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *OneMinuteHandle) Handle() error { //产品品类设备数量统计
	err := stores.WithNoDebug(l.ctx, relationDB.NewProtocolServiceRepo).DownStatus(l.ctx)
	if err != nil {
		l.Error(err)
	}
	return nil
}
