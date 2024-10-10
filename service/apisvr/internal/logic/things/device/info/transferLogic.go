package info

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferLogic {
	return &TransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferLogic) Transfer(req *types.DeviceInfoTransferReq) error {
	_, err := l.svcCtx.DeviceM.DeviceTransfer(l.ctx, utils.Copy[dm.DeviceTransferReq](req))
	return err
}
