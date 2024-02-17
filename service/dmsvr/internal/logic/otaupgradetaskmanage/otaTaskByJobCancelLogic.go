package otaupgradetaskmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskByJobCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaUpgradeTaskRepo
}

func NewOtaTaskByJobCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskByJobCancelLogic {
	return &OtaTaskByJobCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
	}
}

// 取消指定批次下的设备升级作业
func (l *OtaTaskByJobCancelLogic) OtaTaskByJobCancel(in *dm.OTATaskByJobCancelReq) (*dm.Empty, error) {
	var taskStatusList []int
	if in.CancelQueuedTask == 1 {
		taskStatusList = append(taskStatusList, msgOta.UpgradeStatusQueued)
	}
	if in.CancelInProgressTask == 1 {
		taskStatusList = append(taskStatusList, msgOta.UpgradeStatusInProgress)
	}
	if in.CancelNotifiedTask == 1 {
		taskStatusList = append(taskStatusList, msgOta.UpgradeStatusNotified)
	}
	if in.CancelUnconfirmedTask == 1 {
		taskStatusList = append(taskStatusList, msgOta.UpgradeStatusConfirm)
	}
	filter := relationDB.OtaUpgradeTaskFilter{JobId: in.JobId, TaskStatusList: taskStatusList, WithScheduleTime: in.CancelScheduledTask == 1}
	updateData := make(map[string]interface{})
	updateData["task_status"] = msgOta.UpgradeStatusCanceled
	err := l.OtDB.BatchUpdateField(l.ctx, filter, updateData)
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfo BatchUpdate failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Empty{}, nil
}
