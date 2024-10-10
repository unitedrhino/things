package device

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLogic {
	return &CancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLogic) Cancel(req *types.OtaFirmwareDeviceCancelReq) error {
	_, err := l.svcCtx.OtaM.OtaFirmwareDeviceCancel(l.ctx, utils.Copy[dm.OtaFirmwareDeviceCancelReq](req))
	return err
}
