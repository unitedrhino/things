package otaupgradetaskmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaUnfinishedTaskByDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaUpgradeTaskRepo
}

func NewOtaUnfinishedTaskByDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaUnfinishedTaskByDeviceIndexLogic {
	return &OtaUnfinishedTaskByDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
	}
}

// 查询指定设备下，未完成状态的设备升级作业
func (l *OtaUnfinishedTaskByDeviceIndexLogic) OtaUnfinishedTaskByDeviceIndex(in *dm.OTAUnfinishedTaskByDeviceIndexReq) (*dm.OTAUnfinishedTaskByDeviceIndexResp, error) {
	taskStatusList := []int{msgOta.UpgradeStatusConfirm, msgOta.UpgradeStatusInProgress, msgOta.UpgradeStatusQueued, msgOta.UpgradeStatusNotified}
	filter := relationDB.OtaUpgradeTaskFilter{
		ProductId:      in.ProductId,
		TaskStatusList: taskStatusList,
		//ModuleName:     in.ModuleName,
		DeviceName: in.DeviceName,
	}
	var otaUpTaskInfo []*dm.OtaUpTaskInfo
	otaTask, err := l.OtDB.FindByFilter(l.ctx, filter, nil)
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	_ = copier.Copy(&otaUpTaskInfo, &otaTask)
	return &dm.OTAUnfinishedTaskByDeviceIndexResp{OtaUpTaskInfoList: otaUpTaskInfo}, nil
}
