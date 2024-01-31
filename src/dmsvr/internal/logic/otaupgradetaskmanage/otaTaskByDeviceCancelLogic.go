package otaupgradetaskmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskByDeviceCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaUpgradeTaskRepo
}

func NewOtaTaskByDeviceCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskByDeviceCancelLogic {
	return &OtaTaskByDeviceCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
	}
}

// 取消指定OTA升级包下状态为待确认、待推送、已推送、升级中状态的设备升级作业
func (l *OtaTaskByDeviceCancelLogic) OtaTaskByDeviceCancel(in *dm.OTATaskByDeviceCancelReq) (*dm.Response, error) {
	taskStatusList := []int{msgOta.UpgradeStatusConfirm, msgOta.UpgradeStatusInProgress, msgOta.UpgradeStatusQueued, msgOta.UpgradeStatusNotified}
	filter := relationDB.OtaUpgradeTaskFilter{
		JobId:          in.JobId,
		ProductId:      in.ProductId,
		TaskStatusList: taskStatusList,
	}
	updateData := make(map[string]interface{})
	err := l.OtDB.BatchUpdateField(l.ctx, filter, updateData)
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfo BatchUpdate failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
