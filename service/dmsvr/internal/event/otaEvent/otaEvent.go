package otaEvent

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/stores"
	otamanagelogic "github.com/i-Things/things/service/dmsvr/internal/logic/otamanage"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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
		j := job
		ctxs.GoNewCtx(o.ctx, func(ctx context.Context) {
			err := otamanagelogic.NewSendMessageToDevicesLogic(ctx, o.svcCtx).PushMessageToDevices(j)
			if err != nil {
				o.Error(err)
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
