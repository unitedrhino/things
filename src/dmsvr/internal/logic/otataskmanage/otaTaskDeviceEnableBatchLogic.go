package otataskmanagelogic

import (
	"context"

	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceEnableBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceEnableBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceEnableBatchLogic {
	return &OtaTaskDeviceEnableBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取当前可执行批次信息
func (l *OtaTaskDeviceEnableBatchLogic) OtaTaskDeviceEnableBatch(in *dm.OtaTaskBatchReq) (*dm.OtaTaskBatchResp, error) {
	var (
		otaTaskDeviceInfo *relationDB.DmOtaTaskDevices
		otDB              = relationDB.NewOtaTaskDevicesRepo(l.ctx)
	)
	var err error
	if in.ID > 0 {
		otaTaskDeviceInfo, err = otDB.FindOne(l.ctx, in.ID)
	} else {
		otaTaskDeviceInfo, err = otDB.FindEnableBatch(l.ctx, relationDB.OtaTaskDevicesFilter{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
			//Version:    in.GetVersion(),
		})
	}
	if err != nil {
		return nil, err
	}
	compareValue := utils.VersionCompare(otaTaskDeviceInfo.TargetVersion, in.GetVersion())

	if compareValue == 0 {
		//升级成功
		otaTaskDeviceInfo.Status = 501
		otDB.Update(l.ctx, otaTaskDeviceInfo)
		return nil, nil
	}
	if compareValue == -1 {
		// 升级失败，版本较高
		otaTaskDeviceInfo.Status = 601
		otaTaskDeviceInfo.Step = -1
		otaTaskDeviceInfo.Version = in.GetVersion()
		otaTaskDeviceInfo.Desc = "设备版本高于升级包版本"
		otDB.Update(l.ctx, otaTaskDeviceInfo)
		return nil, nil
	}
	return &dm.OtaTaskBatchResp{
		ID:         otaTaskDeviceInfo.ID,
		TaskUid:    otaTaskDeviceInfo.TaskUid,
		FirmwareID: otaTaskDeviceInfo.FirmwareID,
	}, nil
}
