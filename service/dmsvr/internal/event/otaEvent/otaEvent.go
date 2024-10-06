package otaEvent

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	otamanagelogic "gitee.com/i-Things/things/service/dmsvr/internal/logic/otamanage"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
