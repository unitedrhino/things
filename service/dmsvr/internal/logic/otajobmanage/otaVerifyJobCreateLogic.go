package otajobmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaVerifyJobCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB  *relationDB.OtaJobRepo
	OfDB  *relationDB.OtaFirmwareRepo
	OtDB  *relationDB.OtaUpgradeTaskRepo
	DiDB  *relationDB.DeviceInfoRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewOtaVerifyJobCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaVerifyJobCreateLogic {
	return &OtaVerifyJobCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
	}
}

// 验证升级包
func (l *OtaVerifyJobCreateLogic) OtaVerifyJobCreate(in *dm.OtaFirmwareVerifyReq) (*dm.UpgradeJobResp, error) {
	if len(in.DeviceNames) > 10 {
		return nil, errors.OtaDeviceNumError
	}
	var dmOtaJob relationDB.DmOtaJob
	err := copier.Copy(&dmOtaJob, &in)
	if err != nil {
		l.Errorf("%s.Copy OtaFirmwareVerify err=%v", utils.FuncName(), err)
		return nil, err
	}
	dmOtaJob.JobType = msgOta.ValidateUpgrade
	err = l.OjDB.Insert(l.ctx, &dmOtaJob)
	if err != nil {
		l.Errorf("AddVerifyJob.Insert err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	deviceNames := in.DeviceNames
	deviceList, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{DeviceNames: deviceNames}, nil)
	firmware, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareFilter{FirmwareID: in.FirmwareId})
	for _, device := range deviceList {
		dmOtaTask := relationDB.DmOtaUpgradeTask{
			FirmwareId:  in.FirmwareId,
			DeviceName:  device.DeviceName,
			JobId:       dmOtaJob.ID,
			DestVersion: firmware.Version,
			SrcVersion:  device.Version,
			ProductId:   device.ProductID,
			TaskDesc:    firmware.Desc,
			TaskStatus:  msgOta.UpgradeStatusQueued,
		}
		err := l.OtDB.Insert(l.ctx, &dmOtaTask)
		if err != nil {
			l.Errorf("AddVerifyTask.Insert err=%+v", err)
			return nil, errors.System.AddDetail(err)
		}
	}
	//发送消息给设备
	err = NewSendMessageToDevicesLogic(l.ctx, l.svcCtx).PushMessageToDevices(deviceList, firmware)
	if err != nil {
		return nil, err
	}
	return &dm.UpgradeJobResp{JobId: dmOtaJob.ID, UtcCreate: utils.ToYYMMddHHSSByTime(dmOtaJob.CreatedTime)}, nil
}
