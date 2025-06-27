package otamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/things/share/topics"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	if in.UpgradeType == msgOta.DynamicUpgrade && len(in.SrcVersions) == 0 {
		return nil, errors.Parameter.AddMsg("动态升级需要填写待升级的版本")
	}
	if in.UpgradeType == msgOta.DynamicUpgrade { //同一个固件下只能有一个动态升级包
		total, err := relationDB.NewOtaJobRepo(l.ctx).CountByFilter(l.ctx, relationDB.OtaJobFilter{
			FirmwareID:  in.FirmwareID,
			UpgradeType: msgOta.DynamicUpgrade,
			Statues:     []int64{msgOta.JobStatusPlanned, msgOta.JobStatusInProgress},
		})
		if err != nil {
			return nil, err
		}
		if total > 0 {
			return nil, errors.Parameter.AddMsg("一个固件下只能有一个动态升级包")
		}
	}
	var dmOtaJob relationDB.DmOtaFirmwareJob
	err := utils.CopyE(&dmOtaJob, &in)
	if err != nil {
		l.Errorf("%s.CopyE StaticUpgradeJob err=%v", utils.FuncName(), err)
		return nil, err
	}
	dmOtaJob.Status = msgOta.JobStatusInProgress
	if dmOtaJob.UpgradeType == msgOta.StaticUpgrade && dmOtaJob.Static.ScheduleTime != 0 {
		//延时执行
		dmOtaJob.Status = msgOta.JobStatusPlanned
	}
	fi, err := l.OfDB.FindOne(l.ctx, in.FirmwareID)
	if err != nil {
		return nil, err
	}
	if utils.SliceIn(fi.Version, in.SrcVersions...) {
		return nil, errors.Parameter.AddMsg("待升级版本号不能和需要升级的版本号相同")
	}

	dmOtaJob.ProductID = fi.ProductID
	//var  []*dm.StaticUpgradeDeviceInfo
	devicePos, err := l.getDevice(in, fi)
	if err != nil {
		return nil, err
	}
	if dmOtaJob.UpgradeType == msgOta.StaticUpgrade && len(devicePos) == 0 {
		return nil, errors.Parameter.AddMsg("没有符合条件的设备")
	}

	var deviceNames []string
	for _, v := range devicePos {
		deviceNames = append(deviceNames, v.DeviceName)
	}
	var confirmDevices []string
	var clearConfirmDevices []string

	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewOtaJobRepo(tx).Insert(l.ctx, &dmOtaJob)
		if err != nil {
			return err
		}
		if len(devicePos) == 0 {
			return nil
		}
		otDB := relationDB.NewOtaFirmwareDeviceRepo(tx)
		oldDevices, err := otDB.FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			ProductID:   fi.ProductID,
			DeviceNames: deviceNames,
			Statues: []int64{
				msgOta.DeviceStatusConfirm, msgOta.DeviceStatusInProgress, msgOta.DeviceStatusQueued, msgOta.DeviceStatusNotified, msgOta.DeviceStatusFailure},
		}, nil)
		if err != nil {
			return err
		}
		var oldDevicesMap = map[string]*relationDB.DmOtaFirmwareDevice{}
		for _, v := range oldDevices {
			oldDevicesMap[v.DeviceName] = v
		}
		var otaDevices []*relationDB.DmOtaFirmwareDevice
		for _, device := range devicePos {
			status := msgOta.DeviceStatusQueued
			detail := "待推送"
			if in.IsNeedConfirm == def.True {
				status = msgOta.DeviceStatusConfirm
				detail = "待确认"
			}
			od := oldDevicesMap[device.DeviceName]
			if od != nil {
				switch od.Status {
				case msgOta.DeviceStatusInProgress, msgOta.DeviceStatusNotified:
					status = msgOta.DeviceStatusFailure
					detail = "其他任务正在升级中"
				case msgOta.DeviceStatusFailure:
					od.Detail = od.Detail + "-其他任务启动"
					od.Status = msgOta.DeviceStatusCanceled
					err := otDB.Update(l.ctx, od)
					if err != nil {
						return err
					}
					if status == msgOta.DeviceStatusConfirm {
						confirmDevices = append(confirmDevices, device.DeviceName)
					}
				case msgOta.DeviceStatusConfirm, msgOta.DeviceStatusQueued:
					if in.IsOverwriteMode != def.True { //如果是不覆盖则直接失败
						status = msgOta.DeviceStatusFailure
						detail = "其他任务正在等待升级中"
					} else {
						od.Status = msgOta.DeviceStatusCanceled
						od.Detail = "其他任务启动取消该任务"
						err := otDB.Update(l.ctx, od)
						if err != nil {
							return err
						}
						if status == msgOta.DeviceStatusConfirm {
							confirmDevices = append(confirmDevices, device.DeviceName)
						}
					}
				}
			} else if status == msgOta.DeviceStatusConfirm {
				confirmDevices = append(confirmDevices, device.DeviceName)
			}

			if status == msgOta.DeviceStatusQueued { //如果需要执行且不需要确认,则需要将该设备的确认状态清除
				clearConfirmDevices = append(clearConfirmDevices, device.DeviceName)
			}

			otaDevices = append(otaDevices, &relationDB.DmOtaFirmwareDevice{
				FirmwareID:  in.FirmwareID,
				ProductID:   device.ProductID,
				DeviceName:  device.DeviceName,
				JobID:       dmOtaJob.ID,
				SrcVersion:  device.Version,
				DestVersion: fi.Version,
				Status:      status,
				Detail:      detail,
			})
		}
		err = otDB.MultiInsert(l.ctx, otaDevices)
		if err != nil {
			return err
		}
		if len(clearConfirmDevices) > 0 {
			err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID, DeviceNames: clearConfirmDevices},
				map[string]any{"need_confirm_job_id": 0, "need_confirm_version": ""})
			if err != nil {
				return err
			}
		}
		if len(confirmDevices) > 0 {
			err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID, DeviceNames: confirmDevices},
				map[string]any{"need_confirm_job_id": dmOtaJob.ID, "need_confirm_version": fi.Version})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if dmOtaJob.Status == msgOta.JobStatusPlanned {
		_, err := l.svcCtx.TimedM.TaskSend(l.ctx, &timedjob.TaskSendReq{
			GroupCode: def.TimedUnitedRhinoQueueGroupCode,
			Code:      "iThingsDmOtaJobDelayRun",
			Option: &timedjob.TaskSendOption{
				ProcessIn: dmOtaJob.Static.ScheduleTime,
			},
			ParamQueue: &timedjob.TaskParamQueue{
				Topic:   topics.DmOtaJobDelayRun,
				Payload: cast.ToString(dmOtaJob.ID),
			},
		})
		if err != nil {
			l.Error(err)
		}
	}
	if len(confirmDevices) > 0 {
		for _, v := range confirmDevices {
			err := l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{ProductID: fi.ProductID, DeviceName: v}, nil)
			if err != nil {
				l.Error(err)
			}
		}
	}
	if len(clearConfirmDevices) > 0 {
		for _, v := range clearConfirmDevices {
			err := l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{ProductID: fi.ProductID, DeviceName: v}, nil)
			if err != nil {
				l.Error(err)
			}
		}
	}
	return &dm.WithID{Id: dmOtaJob.ID}, err
}

func (l *OtaFirmwareJobCreateLogic) getDevice(in *dm.OtaFirmwareJobInfo, fi *relationDB.DmOtaFirmwareInfo) ([]*relationDB.DmDeviceInfo, error) {
	var devices []*relationDB.DmDeviceInfo
	selection := in.TargetSelection
	switch selection {
	case msgOta.GroupUpgrade:
		ret, err := l.GdDB.FindByFilter(l.ctx, relationDB.GroupDeviceFilter{ProductID: fi.ProductID, Versions: in.SrcVersions, GroupIDs: []int64{cast.ToInt64(in.Target)}, WithDevice: true}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			if v.Device != nil && v.Device.Version != fi.Version {
				devices = append(devices, v.Device)
			}
		}
	case msgOta.AreaUpgrade: //区域升级
		ret, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID, Versions: in.SrcVersions, AreaIDs: []int64{cast.ToInt64(in.Target)}, TenantCodes: in.TenantCodes}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			if v.Version == fi.Version {
				continue
			}
			devices = append(devices, v)
		}
	case msgOta.AllUpgrade, msgOta.GrayUpgrade:
		f := relationDB.DeviceFilter{ProductID: fi.ProductID, Versions: in.SrcVersions, TenantCodes: in.TenantCodes}
		var size int64 = 99999999
		if selection == msgOta.GrayUpgrade {
			total, err := l.DiDB.CountByFilter(l.ctx, f)
			if err != nil {
				return nil, err
			}
			size = int64(float64(total)*(float64(in.Static.GrayPercent)/10000)) + 1
		} else {
			var err error
			size, err = l.DiDB.CountByFilter(l.ctx, f)
			if err != nil {
				return nil, err
			}
		}
		var p int64 = 0
		var psize int64 = 500
		var finish bool
		for p*500 < size && !finish {
			p++
			ret, err := l.DiDB.FindByFilter(l.ctx, f, &stores.PageInfo{
				Page: p,
				Size: psize,
			})
			if err != nil {
				return nil, err
			}
			if selection == msgOta.GrayUpgrade && psize*p > size { //如果是最后一页,获取剩余数量
				s := size - (psize * (p - 1))
				ret = ret[0:s]
			}
			for _, v := range ret {
				if v.Version == fi.Version {
					continue
				}
				devices = append(devices, v)
			}
		}

	case msgOta.SpecificUpgrade:
		ret, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: fi.ProductID, Versions: in.SrcVersions, DeviceNames: in.Static.TargetDeviceNames, TenantCodes: in.TenantCodes}, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range ret {
			if v.Version == fi.Version {
				continue
			}
			devices = append(devices, v)
		}
	default:
		return nil, errors.Parameter.AddMsgf("不支持的升级方式:%v", selection)
	}
	return devices, nil
}

//func (l *OtaFirmwareJobCreateLogic) OtaFirmwareStaticJobCreate(in *dm.OtaFirmwareJobInfo) (*dm.WithID, error) {
//
//
//	return &dm.WithID{Id: dmOtaJob.ID}, nil
//}
//func (l *OtaFirmwareJobCreateLogic) OtaFirmwareDynamicJobCreate(in *dm.OtaFirmwareJobInfo) (*dm.WithID, error) {
//	var dmOtaJob relationDB.DmOtaFirmwareJob
//	err := copier.CopyE(&dmOtaJob, &in)
//	if err != nil {
//		l.Errorf("%s.CopyE DynamicUpgradeJob err=%v", utils.FuncName(), err)
//		return nil, err
//	}
//	dmOtaJob.Type = msgOta.BatchUpgrade
//	dmOtaJob.UpgradeType = msgOta.DynamicUpgrade
//	selection := in.TargetSelection
//	var deviceInfoList []*relationDB.DmDeviceInfo
//	//定向升级
//	if selection == msgOta.SpecificUpgrade {
//		_ = copier.CopyE(&deviceInfoList, &in.DeviceInfos)
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
//			Msg: msgOta.UpgradeStatusQueued,
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
