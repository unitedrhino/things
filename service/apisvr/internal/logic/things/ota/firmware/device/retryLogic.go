package device

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRetryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetryLogic {
	return &RetryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RetryLogic) Retry(req *types.OtaFirmwareDeviceRetryReq) error {
	_, err := l.svcCtx.OtaM.OtaFirmwareDeviceRetry(l.ctx, utils.Copy[dm.OtaFirmwareDeviceRetryReq](req))
	return err
}
