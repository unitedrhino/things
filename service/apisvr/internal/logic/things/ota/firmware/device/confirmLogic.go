package device

import (
	"context"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmLogic {
	return &ConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmLogic) Confirm(req *types.OtaFirmwareDeviceConfirmReq) error {
	_, err := l.svcCtx.OtaM.OtaFirmwareDeviceConfirm(l.ctx, &dm.OtaFirmwareDeviceConfirmReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
	})

	return err
}
