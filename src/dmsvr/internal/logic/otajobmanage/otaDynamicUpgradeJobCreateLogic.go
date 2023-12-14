package otajobmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaDynamicUpgradeJobCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
	GdDB *relationDB.GroupDeviceRepo
	OtDB *relationDB.OtaUpgradeTaskRepo
	OjDB *relationDB.OtaJobRepo
}

func NewOtaDynamicUpgradeJobCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaDynamicUpgradeJobCreateLogic {
	return &OtaDynamicUpgradeJobCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		GdDB:   relationDB.NewGroupDeviceRepo(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// 创建动态升级批次
func (l *OtaDynamicUpgradeJobCreateLogic) OtaDynamicUpgradeJobCreate(in *dm.DynamicUpgradeJobReq) (*dm.UpgradeJobResp, error) {
	var dmOtaJob relationDB.DmOtaJob
	err := copier.Copy(&dmOtaJob, &in)
	if err != nil {
		l.Errorf("%s.Copy DynamicUpgradeJob err=%v", utils.FuncName(), err)
		return nil, err
	}
	dmOtaJob.JobType = msgOta.BatchUpgrade
	dmOtaJob.UpgradeType = msgOta.DynamicUpgrade
	selection := in.TargetSelection

	var deviceInfoList []*relationDB.DmDeviceInfo
	//区域升级
	if selection == msgOta.AreaUpgrade {
		//todo
		//全量升级
	} else if selection == msgOta.AllUpgrade {
		deviceInfoList, err = l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductId}, nil)
		//分组升级
	} else if selection == msgOta.GroupUpgrade {
		gd, err := l.GdDB.FindByFilter(l.ctx, relationDB.GroupDeviceFilter{GroupID: in.GroupId, ProductID: in.ProductId, WithDevice: true}, nil)
		if err != nil {
			l.Errorf("%s.DeviceInfo.GroupDeviceInfoRead failure err=%+v", utils.FuncName(), err)
			return nil, err
		}
		for _, v := range gd {
			deviceInfoList = append(deviceInfoList, v.Device)
		}
	}

	for _, device := range deviceInfoList {
		dmOtaTask := relationDB.DmOtaUpgradeTask{
			FirmwareId: in.FirmwareId,
			DeviceName: device.DeviceName,
			JobId:      dmOtaJob.ID,
			SrcVersion: device.Version,
			ProductId:  device.ProductID,
			TaskStatus: msgOta.UpgradeStatusQueued,
		}
		err := l.OtDB.Insert(l.ctx, &dmOtaTask)
		if err != nil {
			l.Errorf("AddDynamicTask.Insert err=%+v", err)
			return nil, errors.System.AddDetail(err)
		}
	}
	return &dm.UpgradeJobResp{JobId: dmOtaJob.ID, UtcCreate: utils.ToYYMMddHHSSByTime(dmOtaJob.CreatedTime)}, nil

}
