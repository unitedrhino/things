package otajobmanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	server "github.com/i-Things/things/src/dmsvr/internal/server/deviceinteract"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
	firmwareFiles, err := l.OffDB.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{FirmwareID: in.FirmwareId}, nil)
	var files []msgOta.File
	_ = copier.Copy(&files, &firmwareFiles)
	for _, device := range deviceList {
		dmOtaTask := relationDB.DmOtaUpgradeTask{
			FirmwareId: in.FirmwareId,
			DeviceName: device.DeviceName,
			JobId:      dmOtaJob.ID,
			SrcVersion: device.Version,
			ProductId:  device.ProductID,
			TaskDesc:   firmware.Desc,
			TaskStatus: msgOta.UpgradeStatusQueued,
		}
		err := l.OtDB.Insert(l.ctx, &dmOtaTask)
		if err != nil {
			l.Errorf("AddVerifyTask.Insert err=%+v", err)
			return nil, errors.System.AddDetail(err)
		}
	}
	if in.NeedPush == 1 {
		//nats,还有状态要改变
		otaUpgrade := msgOta.Upgrade{
			CommonMsg: deviceMsg.CommonMsg{},
			Params: msgOta.UpgradeParams{
				Version:    firmware.Version,
				IsDiff:     firmware.IsDiff,
				SignMethod: firmware.SignMethod,
				Module:     firmware.Module,
				Files:      files,
			},
		}
		//往nats发送消息
		for _, device := range deviceList {
			topic := fmt.Sprintf("%ota/down/upgrade/%s/%s", in.ProductID, device.DeviceName)
			sendMsg := dm.SendMsgReq{
				Topic:   topic,
				Payload: otaUpgrade.AddStatus(err).Bytes(),
			}
			_, err := server.NewDeviceInteractServer(l.svcCtx).SendMsg(l.ctx, &sendMsg)
			if err != nil {
				l.Errorf("错误是", err)
				logx.Infof("消息发送失败")
				//return nil,err
			}
		}
	}
	return &dm.UpgradeJobResp{JobId: dmOtaJob.ID, UtcCreate: utils.ToYYMMddHHSSByTime(dmOtaJob.CreatedTime)}, nil
}
