package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareJobCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB  *relationDB.OtaJobRepo
	OfDB  *relationDB.OtaFirmwareInfoRepo
	OtDB  *relationDB.OtaFirmwareDeviceRepo
	DiDB  *relationDB.DeviceInfoRepo
	OffDB *relationDB.OtaFirmwareFileRepo
	GdDB  *relationDB.GroupDeviceRepo
}

func NewOtaFirmwareJobCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareJobCreateLogic {
	return &OtaFirmwareJobCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
		OtDB:   relationDB.NewOtaFirmwareDeviceRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareInfoRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
		GdDB:   relationDB.NewGroupDeviceRepo(ctx),
	}
}

// 创建静态升级批次
func (l *OtaFirmwareJobCreateLogic) OtaFirmwareJobCreate(in *dm.OtaFirmwareJobInfo) (*dm.WithID, error) {
	var dmOtaJob relationDB.DmOtaFirmwareJob
	err := copier.Copy(&dmOtaJob, &in)
	if err != nil {
		l.Errorf("%s.Copy StaticUpgradeJob err=%v", utils.FuncName(), err)
		return nil, err
	}

	selection := in.TargetSelection
	fi, err := l.OfDB.FindOne(l.ctx, in.FirmwareID)
	if err != nil {
		return nil, err
	}
	//var  []*dm.StaticUpgradeDeviceInfo
	var devices []*relationDB.DmDeviceInfo
	switch selection {
	case msgOta.AreaUpgrade: //区域升级
	//todo
	case msgOta.AllUpgrade:
		ret, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID}, nil)
		if err != nil {
			return nil, err
		}
		devices = append(devices, ret...)

	case msgOta.SpecificUpgrade:
		ret, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID, DeviceNames: in.Static.TargetDeviceNames}, nil)
		if err != nil {
			return nil, err
		}
		devices = append(devices, ret...)
	case msgOta.GrayUpgrade:
		ret, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID}, nil)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		devices = append(devices, ret...)
	}

	err = l.OjDB.Insert(l.ctx, &dmOtaJob)
	if err != nil {
		return nil, err
	}
	var otaDevices []*relationDB.DmOtaFirmwareDevice
	for _, device := range devices {
		otaDevices = append(otaDevices, &relationDB.DmOtaFirmwareDevice{
			FirmwareID: in.FirmwareID,
			DeviceName: device.DeviceName,
			JobID:      dmOtaJob.ID,
			SrcVersion: device.Version,
			ProductID:  device.ProductID,
			Status:     msgOta.UpgradeStatusQueued,
		})
	}
	err = l.OtDB.MultiInsert(l.ctx, otaDevices)
	if err != nil {
		return nil, err
	}
	return &dm.WithID{Id: dmOtaJob.ID}, nil

}

//func (l *OtaFirmwareJobCreateLogic) OtaFirmwareStaticJobCreate(in *dm.OtaFirmwareJobInfo) (*dm.WithID, error) {
//
//
//	return &dm.WithID{Id: dmOtaJob.ID}, nil
//}
//func (l *OtaFirmwareJobCreateLogic) OtaFirmwareDynamicJobCreate(in *dm.OtaFirmwareJobInfo) (*dm.WithID, error) {
//	var dmOtaJob relationDB.DmOtaFirmwareJob
//	err := copier.Copy(&dmOtaJob, &in)
//	if err != nil {
//		l.Errorf("%s.Copy DynamicUpgradeJob err=%v", utils.FuncName(), err)
//		return nil, err
//	}
//	dmOtaJob.Type = msgOta.BatchUpgrade
//	dmOtaJob.UpgradeType = msgOta.DynamicUpgrade
//	selection := in.TargetSelection
//	var deviceInfoList []*relationDB.DmDeviceInfo
//	//定向升级
//	if selection == msgOta.SpecificUpgrade {
//		_ = copier.Copy(&deviceInfoList, &in.DeviceInfos)
//		//区域升级
//	} else if selection == msgOta.AreaUpgrade {
//		//todo
//		//全量升级
//	} else if selection == msgOta.AllUpgrade {
//		deviceInfoList, err = l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID}, nil)
//		//分组升级
//	} else if selection == msgOta.GroupUpgrade {
//		gd, err := l.GdDB.FindByFilter(l.ctx, relationDB.GroupDeviceFilter{GroupIDs: []int64{in.GroupID}, ProductID: in.ProductID, WithDevice: true}, nil)
//		if err != nil {
//			l.Errorf("%s.DeviceInfo.GroupDeviceInfoRead failure err=%+v", utils.FuncName(), err)
//			return nil, err
//		}
//		for _, v := range gd {
//			deviceInfoList = append(deviceInfoList, v.Device)
//		}
//	}
//	for _, device := range deviceInfoList {
//		dmOtaTask := relationDB.DmOtaFirmwareDevice{
//			FirmwareID: in.FirmwareID,
//			DeviceName: device.DeviceName,
//			JobID:      dmOtaJob.ID,
//			SrcVersion: device.Version,
//			ProductID:  device.ProductID,
//			Status: msgOta.UpgradeStatusQueued,
//		}
//		err := l.OtDB.Insert(l.ctx, &dmOtaTask)
//		if err != nil {
//			l.Errorf("AddDynamicTask.Insert err=%+v", err)
//			return nil, errors.System.AddDetail(err)
//		}
//	}
//	//发送消息给设备
//	firmware, err := l.OfDB.FindOne(l.ctx, in.FirmwareID)
//	err = NewSendMessageToDevicesLogic(l.ctx, l.svcCtx).PushMessageToDevices(deviceInfoList, firmware)
//	if err != nil {
//		return nil, err
//	}
//	return &dm.WithID{Id: dmOtaJob.ID}, nil
//}
