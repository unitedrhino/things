package otamanagelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareDeviceRetryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaFirmwareDeviceRetryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareDeviceRetryLogic {
	return &OtaFirmwareDeviceRetryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重新升级指定批次下升级失败或升级取消的设备升级作业
func (l *OtaFirmwareDeviceRetryLogic) OtaFirmwareDeviceRetry(in *dm.OtaFirmwareDeviceRetryReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
