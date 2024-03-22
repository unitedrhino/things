package otaupgradetaskmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaUnfinishedTaskByDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaFirmwareDeviceRepo
}

func NewOtaUnfinishedTaskByDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaUnfinishedTaskByDeviceIndexLogic {
	return &OtaUnfinishedTaskByDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaFirmwareDeviceRepo(ctx),
	}
}

//
//// 查询指定设备下，未完成状态的设备升级作业
//func (l *OtaUnfinishedTaskByDeviceIndexLogic) OtaUnfinishedTaskByDeviceIndex(in *dm.OTAUnfinishedTaskByDeviceIndexReq) (*dm.OTAUnfinishedTaskByDeviceIndexResp, error) {
//	taskStatusList := []int{msgOta.UpgradeStatusConfirm, msgOta.UpgradeStatusInProgress, msgOta.UpgradeStatusQueued, msgOta.UpgradeStatusNotified}
//	filter := relationDB.OtaFirmwareDeviceFilter{
//		ProductID:      in.ProductID,
//		Statues: taskStatusList,
//		//ModuleName:     in.ModuleName,
//		DeviceName: in.DeviceName,
//	}
//	var otaUpTaskInfo []*dm.OtaUpTaskInfo
//	otaTask, err := l.OtDB.FindByFilter(l.ctx, filter, nil)
//	if err != nil {
//		l.Errorf("%s.TaskInfo.TaskInfo failure err=%+v", utils.FuncName(), err)
//		return nil, err
//	}
//	_ = copier.Copy(&otaUpTaskInfo, &otaTask)
//	return &dm.OTAUnfinishedTaskByDeviceIndexResp{OtaUpTaskInfoList: otaUpTaskInfo}, nil
//}
