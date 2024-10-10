package otamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	if di.NeedConfirmVersion == "" {
		return nil, errors.OtaCancleStatusError.WithMsg("已经升级完成")
	}
	f := relationDB.OtaFirmwareDeviceFilter{
		ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}, DestVersion: di.NeedConfirmVersion, JobID: di.NeedConfirmJobID}
	dev, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindOneByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	if utils.SliceIn(dev.Status, msgOta.DeviceStatusNotified, msgOta.DeviceStatusInProgress) {
		return nil, errors.OtaRetryStatusError.WithMsg("正在升级中,请耐心等待")
	}
	if di.Version.GetValue() == dev.DestVersion {
		return nil, errors.OtaCancleStatusError.WithMsg("已经升级成功")
	}
	dev.Status = msgOta.DeviceStatusQueued
	dev.Detail = "手动执行待推送"
	err = relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, dev)
	if err != nil {
		return nil, err
	}
	//err = relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}}, map[string]any{
	//	"need_confirm_job_id":  0,
	//	"need_confirm_version": "",
	//})
	return &dm.Empty{}, err
}
