package userdevicelogic

import (
	"context"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceTransferLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceTransferLogic {
	return &UserDeviceTransferLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeviceTransferLogic) UserDeviceTransfer(in *dm.DeviceTransferReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
