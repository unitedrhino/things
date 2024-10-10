package otamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareJobUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB *relationDB.OtaJobRepo
}

func NewOtaFirmwareJobUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareJobUpdateLogic {
	return &OtaFirmwareJobUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// 取消动态升级策略
func (l *OtaFirmwareJobUpdateLogic) OtaFirmwareJobUpdate(in *dm.OtaFirmwareJobInfo) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	otaJob, err := l.OjDB.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	if in.Status == msgOta.JobStatusCanceled && otaJob.Status != in.Status {
		otaJob.Status = in.Status
		err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			err := relationDB.NewOtaFirmwareDeviceRepo(tx).UpdateStatusByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
				FirmwareID: otaJob.FirmwareID,
				JobID:      otaJob.ID,
				Statues:    []int64{msgOta.DeviceStatusFailure, msgOta.DeviceStatusQueued, msgOta.DeviceStatusConfirm}, //需要重试的设备更换为待推送
			}, msgOta.DeviceStatusCanceled, "任务取消,取消升级")
			if err != nil {
				return err
			}
			err = relationDB.NewDeviceInfoRepo(tx).UpdateWithField(l.ctx, relationDB.DeviceFilter{NeedConfirmJobID: otaJob.ID}, map[string]any{
				"need_confirm_job_id":  0,
				"need_confirm_version": "",
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	if in.MaximumPerMinute != 0 {
		otaJob.MaximumPerMinute = in.MaximumPerMinute
	}
	err = l.OjDB.Update(l.ctx, otaJob)
	return &dm.Empty{}, err
}
