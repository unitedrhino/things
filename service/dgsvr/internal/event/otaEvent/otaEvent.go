package otaEvent

import (
	"context"
	"gitee.com/unitedrhino/things/service/dgsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaEvent struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	ctx context.Context
}

func NewOtaEvent(svcCtx *svc.ServiceContext, ctx context.Context) *OtaEvent {
	return &OtaEvent{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (o *OtaEvent) DeviceUpgradePush() error {
	return nil
}

func (o *OtaEvent) JobDelayRun(jobID int64) error {
	o.Info(jobID)
	return nil
}
