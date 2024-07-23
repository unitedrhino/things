package otamanagelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type SendMessageToDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB  *relationDB.DeviceInfoRepo
	GdDB  *relationDB.GroupDeviceRepo
	OtDB  *relationDB.OtaFirmwareDeviceRepo
	OjDB  *relationDB.OtaJobRepo
	OfDB  *relationDB.OtaFirmwareInfoRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewSendMessageToDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageToDevicesLogic {
	return &SendMessageToDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   stores.WithNoDebug(ctx, relationDB.NewOtaJobRepo),
		OtDB:   stores.WithNoDebug(ctx, relationDB.NewOtaFirmwareDeviceRepo),
		DiDB:   stores.WithNoDebug(ctx, relationDB.NewDeviceInfoRepo),
		OfDB:   stores.WithNoDebug(ctx, relationDB.NewOtaFirmwareInfoRepo),
		OffDB:  stores.WithNoDebug(ctx, relationDB.NewOtaFirmwareFileRepo),
	}
}

func (l *SendMessageToDevicesLogic) DevicesTimeout(jobInfo *relationDB.DmOtaFirmwareJob) error {
	firmware := jobInfo.Firmware
	if jobInfo.IsNeedPush != def.True { //只有需要推送的才推送
		return nil
	}
	pushDevice := func(devs []devices.Core, status int64, detail string) {
		for _, df := range devs {
			appMsg := application.OtaReport{
				Device:    df,
				Timestamp: time.Now().UnixMilli(), Status: status, Detail: detail,
			}
			di, err := l.svcCtx.DeviceCache.GetData(l.ctx, df)
			if err != nil {
				l.Error(err)
				continue
			}
			err = l.svcCtx.UserSubscribe.Publish(l.ctx, def.UserSubscribeDeviceOtaReport, appMsg, map[string]any{
				"productID":  di.ProductID,
				"deviceName": di.DeviceName,
			}, map[string]any{
				"projectID": di.ProjectID,
			}, map[string]any{
				"projectID": cast.ToString(di.ProjectID),
				"areaID":    cast.ToString(di.AreaID),
			})
			if err != nil {
				l.Error(err)
			}
		}
	}
	{ //处理超时设备,置为失败
		f := relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: jobInfo.FirmwareID,
			JobID:      jobInfo.ID,
			ProductID:  firmware.ProductID,
			PushTime:   stores.CmpLte(time.Now().Add(-time.Duration(jobInfo.TimeoutInMinutes) * time.Minute)),
			Statues:    []int64{msgOta.DeviceStatusNotified, msgOta.DeviceStatusInProgress}, //只处理待推送的设备
		}
		var devs []devices.Core
		err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			ofdr := relationDB.NewOtaFirmwareDeviceRepo(tx)
			pos, err := ofdr.FindByFilter(l.ctx, f, nil)
			if err != nil {
				l.Error(err)
				return err
			}
			for _, po := range pos {
				devs = append(devs, devices.Core{ProductID: po.ProductID, DeviceName: po.DeviceName})
			}
			err = ofdr.UpdateStatusByFilter(l.ctx, f, msgOta.DeviceStatusFailure, "设备超时") //如果超过了超时时间,则修改为失败
			if err != nil {
				l.Error(err)
			}
			return err
		})
		if err != nil {
			l.Error(err)
		} else {
			pushDevice(devs, msgOta.DeviceStatusFailure, "设备超时")
		}

	}

	if jobInfo.RetryCount > 0 { //处理重试设备
		var devs []devices.Core
		f := relationDB.OtaFirmwareDeviceFilter{
			FirmwareID:      jobInfo.FirmwareID,
			JobID:           jobInfo.ID,
			ProductID:       firmware.ProductID,
			LastFailureTime: stores.CmpLte(time.Now().Add(-time.Minute * time.Duration(jobInfo.RetryInterval))), //失败间隔
			RetryCount:      stores.CmpLt(jobInfo.RetryCount),                                                   //重试次数
			Statues:         []int64{msgOta.DeviceStatusFailure},                                                //需要重试的设备更换为待推送
		}
		err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			ofdr := relationDB.NewOtaFirmwareDeviceRepo(tx)
			pos, err := ofdr.FindByFilter(l.ctx, f, nil)
			if err != nil {
				return err
			}
			for _, po := range pos {
				devs = append(devs, devices.Core{ProductID: po.ProductID, DeviceName: po.DeviceName})
			}
			err = ofdr.UpdateStatusByFilter(l.ctx, f, msgOta.DeviceStatusQueued, "重试推送") //如果超过了超时时间,则修改为失败
			return err
		})
		if err != nil {
			l.Error(err)
		} else {
			pushDevice(devs, msgOta.DeviceStatusQueued, "重试推送")
		}
	}
	{
		var devs []devices.Core
		f := relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: jobInfo.FirmwareID,
			JobID:      jobInfo.ID,
			ProductID:  firmware.ProductID,
			RetryCount: stores.CmpGte(jobInfo.RetryCount),   //重试次数
			Statues:    []int64{msgOta.DeviceStatusFailure}, //需要重试的设备更换为待推送
		}
		err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			ofdr := relationDB.NewOtaFirmwareDeviceRepo(tx)
			pos, err := ofdr.FindByFilter(l.ctx, f, nil)
			if err != nil {
				return err
			}
			for _, po := range pos {
				devs = append(devs, devices.Core{ProductID: po.ProductID, DeviceName: po.DeviceName})
			}
			err = ofdr.UpdateStatusByFilter(l.ctx, f, msgOta.DeviceStatusCanceled, "超过重试次数,取消升级") //如果超过了超时时间,则修改为失败
			return err
		})
		if err != nil {
			l.Error(err)
		} else {
			pushDevice(devs, msgOta.DeviceStatusCanceled, "超过重试次数,取消升级")
		}
	}
	func() {
		total, err := stores.WithNoDebug(l.ctx, relationDB.NewOtaFirmwareDeviceRepo).CountByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: jobInfo.FirmwareID,
			JobID:      jobInfo.ID,
		})
		if err != nil {
			l.Error(err)
			return
		}
		finished, err := stores.WithNoDebug(l.ctx, relationDB.NewOtaFirmwareDeviceRepo).CountByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: jobInfo.FirmwareID,
			JobID:      jobInfo.ID,
			Statues:    []int64{msgOta.DeviceStatusCanceled, msgOta.DeviceStatusSuccess},
		})
		if err != nil {
			l.Error(err)
			return
		}
		if total == finished { //任务完成
			newJob := *jobInfo
			newJob.Status = msgOta.JobStatusCompleted
			err = stores.WithNoDebug(l.ctx, relationDB.NewOtaJobRepo).Update(l.ctx, &newJob)
			if err != nil {
				l.Error(err)
				return
			}
		}
	}()

	return nil
}
func (l *SendMessageToDevicesLogic) PushMessageToDevices(jobInfo *relationDB.DmOtaFirmwareJob) error {
	err := l.DevicesTimeout(jobInfo)
	if err != nil {
		l.Error(err)
	}
	firmware := jobInfo.Firmware
	if jobInfo.IsNeedPush != def.True { //只有需要推送的才推送
		return nil
	}

	deviceList, err := stores.WithNoDebug(l.ctx, relationDB.NewOtaFirmwareDeviceRepo).FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		FirmwareID: jobInfo.FirmwareID,
		JobID:      jobInfo.ID,
		ProductID:  firmware.ProductID,
		IsOnline:   def.True,                           //只有在线的设备才推送
		Statues:    []int64{msgOta.DeviceStatusQueued}, //只处理待推送的设备
	}, &stores.PageInfo{
		Page: 1,
		Size: jobInfo.MaximumPerMinute/(60/5) + 1, //任务5秒钟推送一次
	})
	if err != nil {
		return err
	}
	if len(deviceList) == 0 {
		//todo 这里需要结束任务,没有需要执行的了
		return nil
	}
	firmwareFiles, err := l.OffDB.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{FirmwareID: jobInfo.FirmwareID}, nil)
	if err != nil {
		return err
	}
	data, err := GenUpgradeParams(l.ctx, l.svcCtx, firmware, firmwareFiles)
	if err != nil {
		return err
	}
	MsgToken := devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)
	upgradeMsg := deviceMsg.CommonMsg{
		MsgToken:  MsgToken,
		Method:    msgOta.TypeUpgrade,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
	payload, _ := json.Marshal(upgradeMsg)
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, firmware.ProductID)
	if err != nil {
		return err
	}
	for _, device := range deviceList {
		reqMsg := deviceMsg.PublishMsg{
			Handle:       devices.Ota,
			Type:         msgOta.TypeUpgrade,
			Payload:      payload,
			Timestamp:    time.Now().UnixMilli(),
			ProductID:    device.ProductID,
			DeviceName:   device.DeviceName,
			ProtocolCode: pi.ProtocolCode,
		}
		err := l.svcCtx.PubDev.PublishToDev(l.ctx, &reqMsg)
		if err != nil {
			return err
		}
		device.Status = msgOta.DeviceStatusNotified
		device.Detail = "主动推送"
		device.PushTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err = relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, device)
		if err != nil {
			return err
		}
		core := devices.Core{
			ProductID:  device.ProductID,
			DeviceName: device.DeviceName,
		}
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, core)
		if err != nil {
			l.Error(err)
			return nil
		}
		appMsg := application.OtaReport{
			Device:    core,
			Timestamp: time.Now().UnixMilli(),
			Status:    device.Status,
			Detail:    device.Detail,
			Step:      device.Step,
		}
		err = l.svcCtx.UserSubscribe.Publish(l.ctx, def.UserSubscribeDeviceOtaReport, appMsg, map[string]any{
			"productID":  di.ProductID,
			"deviceName": di.DeviceName,
		}, map[string]any{
			"projectID": di.ProjectID,
		}, map[string]any{
			"projectID": cast.ToString(di.ProjectID),
			"areaID":    cast.ToString(di.AreaID),
		})
	}
	return nil
}

func GenUpgradeParams(ctx context.Context, svcCtx *svc.ServiceContext, firmware *relationDB.DmOtaFirmwareInfo, files []*relationDB.DmOtaFirmwareFile) (*msgOta.UpgradeData, error) {
	if len(files) == 0 {
		return nil, errors.System.AddDetail("升级包下没有文件")
	}
	if len(files) == 1 { //单文件模式
		url, err := svcCtx.OssClient.PublicBucket().GetUrl(files[0].FilePath, true)
		if err != nil {
			return nil, err
		}
		data := msgOta.UpgradeData{
			Version:    firmware.Version,
			IsDiff:     firmware.IsDiff,
			SignMethod: firmware.SignMethod,
			Extra:      firmware.Extra,
			File: &msgOta.File{
				Size:      files[0].Size,
				Name:      files[0].Name,
				FileUrl:   url,
				FileMd5:   files[0].FileMd5,
				Signature: files[0].Signature,
			},
		}
		return &data, nil
	}
	var data = msgOta.UpgradeData{
		Version:    firmware.Version,
		IsDiff:     firmware.IsDiff,
		SignMethod: firmware.SignMethod,
		Extra:      firmware.Extra,
	}
	for _, f := range files {
		url, err := svcCtx.OssClient.PublicBucket().GetUrl(f.FilePath, true)
		if err != nil {
			return nil, err
		}
		data.Files = append(data.Files, &msgOta.File{
			Size:      f.Size,
			Name:      f.Name,
			FileUrl:   url,
			FileMd5:   f.FileMd5,
			Signature: f.Signature,
		})
	}
	return &data, nil
}
