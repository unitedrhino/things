package otamanagelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareDeviceCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaFirmwareDeviceCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareDeviceCancelLogic {
	return &OtaFirmwareDeviceCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消指定批次下的设备升级作业
func (l *OtaFirmwareDeviceCancelLogic) OtaFirmwareDeviceCancel(in *dm.OtaFirmwareDeviceCancelReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
