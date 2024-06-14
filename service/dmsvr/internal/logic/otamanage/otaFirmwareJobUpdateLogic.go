package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

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
	otaJob, err := l.OjDB.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	if in.Status == msgOta.JobStatusCanceled && otaJob.Status != in.Status {
		otaJob.Status = in.Status
		err := stores.WithNoDebug(l.ctx, relationDB.NewOtaFirmwareDeviceRepo).UpdateStatusByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			FirmwareID: otaJob.FirmwareID,
			JobID:      otaJob.ID,
			Statues:    []int64{msgOta.DeviceStatusFailure, msgOta.DeviceStatusConfirm}, //需要重试的设备更换为待推送
		}, msgOta.DeviceStatusCanceled, "任务取消,取消升级") //如果超过了超时时间,则修改为失败
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
