package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

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
	if len(in.DeviceNames) == 0 {
		return nil, errors.Parameter.AddMsg("设备名列表必填")
	}
	err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).BatchUpdateField(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		FirmwareID:  in.FirmwareID,
		JobID:       in.JobID,
		DeviceNames: in.DeviceNames,
		Statues:     []int64{msgOta.DeviceStatusConfirm, msgOta.DeviceStatusQueued},
	}, map[string]interface{}{"status": msgOta.DeviceStatusCanceled})
	return &dm.Empty{}, err
}
