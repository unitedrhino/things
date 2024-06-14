package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareDeviceConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaFirmwareDeviceConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareDeviceConfirmLogic {
	return &OtaFirmwareDeviceConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// app确认设备升级
func (l *OtaFirmwareDeviceConfirmLogic) OtaFirmwareDeviceConfirm(in *dm.OtaFirmwareDeviceConfirmReq) (*dm.Empty, error) {
	f := relationDB.OtaFirmwareDeviceFilter{
		ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}, Statues: []int64{msgOta.DeviceStatusConfirm}}
	dev, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindOneByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	dev.Status = msgOta.DeviceStatusQueued
	err = relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, dev)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}}, map[string]any{
		"need_confirm_job_id":  0,
		"need_confirm_version": "",
	})
	return &dm.Empty{}, err
}