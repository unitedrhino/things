package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	if len(in.DeviceNames) == 0 {
		return nil, errors.Parameter.AddMsg("设备名列表必填")
	}
	err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).BatchUpdateField(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		FirmwareID:  in.FirmwareID,
		JobID:       in.JobID,
		DeviceNames: in.DeviceNames,
		Statues:     []int64{msgOta.DeviceStatusSuccess, msgOta.DeviceStatusFailure, msgOta.DeviceStatusCanceled},
	}, map[string]interface{}{"status": msgOta.DeviceStatusQueued})
	return &dm.Empty{}, err
}
