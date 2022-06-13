package logic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageFirmwareLogic {
	return &ManageFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理产品的固件
func (l *ManageFirmwareLogic) ManageFirmware(in *dm.ManageFirmwareReq) (*dm.FirmwareInfo, error) {
	// todo: add your logic here and delete this line

	return &dm.FirmwareInfo{}, nil
}
