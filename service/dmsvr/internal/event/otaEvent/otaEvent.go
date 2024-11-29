package otaEvent

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/application"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	otamanagelogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/otamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"sync"
	"time"
)

type OtaEvent struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	ctx context.Context
}

func NewOtaEvent(svcCtx *svc.ServiceContext, ctx context.Context) *OtaEvent {
	return &OtaEvent{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (o *OtaEvent) DeviceUpgradePush() error {
	jobs, err := stores.WithNoDebug(o.ctx, relationDB.NewOtaJobRepo).FindByFilter(o.ctx, relationDB.OtaJobFilter{
		WithFirmware: true,
		Statues:      []int64{msgOta.JobStatusInProgress},
		WithFiles:    true,
	}, nil)
	if err != nil {
		return err
	}
	err = o.DevicesTimeout(jobs)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		jj := job
		if job.Firmware == nil { //任务的固件已经被删除了,需要删除该任务及对应的设备
			ctxs.GoNewCtx(o.ctx, func(ctx context.Context) {
				err := stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
					err := relationDB.NewOtaFirmwareDeviceRepo(tx).DeleteByFilter(ctx, relationDB.OtaFirmwareDeviceFilter{
						JobID: jj.ID,
					})
					if err != nil {
						return err
					}
					err = relationDB.NewOtaJobRepo(tx).Delete(ctx, jj.ID)
					return err
				})
				if err != nil {
					logx.WithContext(ctx).Errorf("Device upgrade push err:%+v", err)
				}
			})
			continue
		}
		ctxs.GoNewCtx(o.ctx, func(ctx context.Context) {
			start := time.Now()
			defer func() {
				end := time.Now()
				if end.Sub(start).Seconds() > 2 {
					logx.WithContext(ctx).Slowf("PushMessageToDevices use:%v  job:%v", end.Sub(start), utils.Fmt(jj))
				}
			}()
			err := otamanagelogic.NewSendMessageToDevicesLogic(ctx, o.svcCtx).PushMessageToDevices(jj)
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				logx.WithContext(ctx).Error(err)
			}
		})
	}
	return nil
}

func (l *OtaEvent) DevicesTimeout(jobInfos []*relationDB.DmOtaFirmwareJob) error {
	var wait sync.WaitGroup
	defer func() {
		wait.Wait()
	}()
	pushDevice := func(devs []devices.Core, status int64, detail string) {
		wait.Add(1)
		utils.Go(l.ctx, func() {
			wait.Done()
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
		})
	}

	var (
		expFail           []relationDB.OtaFirmwareDeviceFilter
		retry             []relationDB.OtaFirmwareDeviceFilter
		retryNeedConfirm  []relationDB.OtaFirmwareDeviceFilter
		cancel            []relationDB.OtaFirmwareDeviceFilter
		dynamicUpgradeJob []*relationDB.DmOtaFirmwareJob
	)

	for _, jobInfo := range jobInfos {
		firmware := jobInfo.Firmware
		if jobInfo.Firmware == nil { //任务的固件已经被删除了,需要删除该任务及对应的设备
			ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
				err := stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
					err := relationDB.NewOtaFirmwareDeviceRepo(tx).DeleteByFilter(ctx, relationDB.OtaFirmwareDeviceFilter{
						JobID: jobInfo.ID,
					})
					if err != nil {
						return err
					}
					err = relationDB.NewOtaJobRepo(tx).Delete(ctx, jobInfo.ID)
					return err
				})
				if err != nil {
					logx.WithContext(ctx).Errorf("Device upgrade push err:%+v", err)
				}
			})
			continue
		}

		expFail = append(expFail, relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: jobInfo.FirmwareID,
			JobID:      jobInfo.ID,
			ProductID:  firmware.ProductID,
			PushTime:   stores.CmpLte(time.Now().Add(-time.Duration(jobInfo.TimeoutInMinutes) * time.Minute)),
			Statues:    []int64{msgOta.DeviceStatusNotified, msgOta.DeviceStatusInProgress}, //只处理待推送的设备
		})
		if jobInfo.IsNeedPush != def.True { //只有需要推送的才推送
			return nil
		}

		if jobInfo.RetryCount > 0 { //处理重试设备
			f := relationDB.OtaFirmwareDeviceFilter{
				FirmwareID:      jobInfo.FirmwareID,
				JobID:           jobInfo.ID,
				ProductID:       firmware.ProductID,
				LastFailureTime: stores.CmpLte(time.Now().Add(-time.Minute * time.Duration(jobInfo.RetryInterval))), //失败间隔
				RetryCount:      stores.CmpLt(jobInfo.RetryCount),                                                   //重试次数
				Statues:         []int64{msgOta.DeviceStatusFailure},                                                //需要重试的设备更换为待推送
			}
			if jobInfo.IsNeedConfirm == def.True {
				retryNeedConfirm = append(retryNeedConfirm, f)
			} else {
				retry = append(retry, f)
			}
		}
		cancel = append(cancel, relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: jobInfo.FirmwareID,
			JobID:      jobInfo.ID,
			ProductID:  firmware.ProductID,
			RetryCount: stores.CmpGte(jobInfo.RetryCount),   //重试次数
			Statues:    []int64{msgOta.DeviceStatusFailure}, //需要重试的设备更换为待推送
		})
		if jobInfo.UpgradeType == msgOta.DynamicUpgrade { //动态的需要将后面符合升级标准的加进去
			if time.Now().Second() < 5 { //一分钟执行一次
				dynamicUpgradeJob = append(dynamicUpgradeJob, jobInfo)

			}
		}
		func() { //完成任务
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
	}

	if len(expFail) > 0 { //处理超时设备,置为失败
		var pos []*relationDB.DmOtaFirmwareDevice
		var err error
		err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			ofdr := relationDB.NewOtaFirmwareDeviceRepo(tx)
			pos, err = ofdr.FindByFilters(l.ctx, expFail, nil)
			if err != nil {
				l.Error(err)
				return err
			}
			if len(pos) > 0 {
				err = ofdr.UpdateStatusByFilters(l.ctx, expFail, msgOta.DeviceStatusFailure, "设备超时") //如果超过了超时时间,则修改为失败
				if err != nil {
					l.Error(err)
				}
			}
			return err
		})
		if err != nil {
			l.Error(err)
		} else {
			var devs []devices.Core
			for _, po := range pos {
				devs = append(devs, devices.Core{ProductID: po.ProductID, DeviceName: po.DeviceName})
			}
			pushDevice(devs, msgOta.DeviceStatusFailure, "设备超时")
		}

	}
	{
		status := msgOta.DeviceStatusQueued
		detail := "重试推送"
		handleRetry := func(f []relationDB.OtaFirmwareDeviceFilter) {
			var devs []devices.Core

			var pos []*relationDB.DmOtaFirmwareDevice
			var err error
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				ofdr := relationDB.NewOtaFirmwareDeviceRepo(tx)
				pos, err = ofdr.FindByFilters(l.ctx, f, nil)
				if err != nil {
					return err
				}
				if len(pos) > 0 {
					err = ofdr.UpdateStatusByFilters(l.ctx, f, status, detail) //如果超过了超时时间,则修改为失败
				}
				return err
			})
			if err != nil {
				l.Error(err)
			} else if status == msgOta.DeviceStatusQueued {
				for _, po := range pos {
					devs = append(devs, devices.Core{ProductID: po.ProductID, DeviceName: po.DeviceName})
				}
				pushDevice(devs, status, detail)
			}
		}
		if len(retry) > 0 {
			handleRetry(retry)
		}
		if len(retryNeedConfirm) > 0 {
			status = msgOta.DeviceStatusConfirm
			detail = "升级失败,再次升级等待确认"
			handleRetry(retryNeedConfirm)
		}
	}

	if len(cancel) > 0 {
		var pos []*relationDB.DmOtaFirmwareDevice
		var err error
		err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			ofdr := relationDB.NewOtaFirmwareDeviceRepo(tx)
			pos, err = ofdr.FindByFilters(l.ctx, cancel, nil)
			if err != nil {
				return err
			}
			if len(pos) > 0 {
				err = ofdr.UpdateStatusByFilters(l.ctx, cancel, msgOta.DeviceStatusCanceled, "超过重试次数,取消升级") //如果超过了超时时间,则修改为失败
			}
			return err
		})
		if err != nil {
			l.Error(err)
		} else {
			var devs []devices.Core
			for _, po := range pos {
				devs = append(devs, devices.Core{ProductID: po.ProductID, DeviceName: po.DeviceName})
			}
			pushDevice(devs, msgOta.DeviceStatusCanceled, "超过重试次数,取消升级")
		}
	}
	if len(dynamicUpgradeJob) > 0 {
		for _, jobInfo := range dynamicUpgradeJob {
			err := l.AddDevice(jobInfo)
			if err != nil {
				l.Error(err)
			}
		}
	}

	return nil
}

func (l *OtaEvent) AddDevice(dmOtaJob *relationDB.DmOtaFirmwareJob) error {
	devicePos, err := stores.WithNoDebug(l.ctx, relationDB.NewDeviceInfoRepo).FindByFilter(l.ctx,
		relationDB.DeviceFilter{NotOtaJobID: dmOtaJob.ID, ProductID: dmOtaJob.ProductID, Versions: dmOtaJob.SrcVersions}, nil)
	if err != nil {
		l.Error(err)
	}
	if len(devicePos) == 0 {
		return nil
	}

	var deviceNames []string
	for _, v := range devicePos {
		deviceNames = append(deviceNames, v.DeviceName)
	}
	var confirmDevices []*devices.Core
	var clearConfirmDevices []*devices.Core

	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		otDB := relationDB.NewOtaFirmwareDeviceRepo(tx)
		oldDevices, err := otDB.FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			ProductID:   dmOtaJob.ProductID,
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
			if dmOtaJob.IsNeedConfirm == def.True {
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
						confirmDevices = append(confirmDevices, &devices.Core{
							ProductID:  device.ProductID,
							DeviceName: device.DeviceName,
						})
					}
				case msgOta.DeviceStatusConfirm, msgOta.DeviceStatusQueued:
					if dmOtaJob.IsOverwriteMode != def.True { //如果是不覆盖则直接失败
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
							confirmDevices = append(confirmDevices, &devices.Core{
								ProductID:  device.ProductID,
								DeviceName: device.DeviceName,
							})
						}
					}
				}
			} else if status == msgOta.DeviceStatusConfirm {
				confirmDevices = append(confirmDevices, &devices.Core{
					ProductID:  device.ProductID,
					DeviceName: device.DeviceName,
				})
			}

			if status == msgOta.DeviceStatusQueued { //如果需要执行且不需要确认,则需要将该设备的确认状态清除
				clearConfirmDevices = append(clearConfirmDevices, &devices.Core{
					ProductID:  device.ProductID,
					DeviceName: device.DeviceName,
				})
			}

			otaDevices = append(otaDevices, &relationDB.DmOtaFirmwareDevice{
				FirmwareID:  dmOtaJob.FirmwareID,
				ProductID:   device.ProductID,
				DeviceName:  device.DeviceName,
				JobID:       dmOtaJob.ID,
				SrcVersion:  device.Version,
				DestVersion: dmOtaJob.Firmware.Version,
				Status:      status,
				Detail:      detail,
			})
		}
		err = otDB.MultiInsert(l.ctx, otaDevices)
		if err != nil {
			return err
		}
		if len(clearConfirmDevices) > 0 {
			err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(l.ctx, relationDB.DeviceFilter{Cores: clearConfirmDevices},
				map[string]any{"need_confirm_job_id": 0, "need_confirm_version": ""})
			if err != nil {
				return err
			}
		}
		if len(confirmDevices) > 0 {
			err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(l.ctx, relationDB.DeviceFilter{Cores: confirmDevices},
				map[string]any{"need_confirm_job_id": dmOtaJob.ID, "need_confirm_version": dmOtaJob.Firmware.Version})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if len(confirmDevices) > 0 {
		for _, v := range confirmDevices {
			err := l.svcCtx.DeviceCache.SetData(l.ctx, *v, nil)
			if err != nil {
				l.Error(err)
			}
		}
	}
	if len(clearConfirmDevices) > 0 {
		for _, v := range clearConfirmDevices {
			err := l.svcCtx.DeviceCache.SetData(l.ctx, *v, nil)
			if err != nil {
				l.Error(err)
			}
		}
	}
	return nil
}

func (o *OtaEvent) JobDelayRun(jobID int64) error {
	o.Info(jobID)
	oj, err := relationDB.NewOtaJobRepo(o.ctx).FindOne(o.ctx, jobID)
	if err != nil {
		return err
	}
	oj.Status = msgOta.JobStatusInProgress
	err = relationDB.NewOtaJobRepo(o.ctx).Update(o.ctx, oj)
	return err
}
