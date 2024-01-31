package otaupgradetaskmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskByDeviceNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaUpgradeTaskRepo
}

func NewOtaTaskByDeviceNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskByDeviceNameLogic {
	return &OtaTaskByDeviceNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
	}
}

// 根据设备的name和productId，查询task
func (l *OtaTaskByDeviceNameLogic) OtaTaskByDeviceName(in *dm.OTATaskByDeviceNameReq) (*dm.OtaUpTaskInfo, error) {
	filter := relationDB.OtaUpgradeTaskFilter{
		ProductId:  in.ProductId,
		DeviceName: in.DeviceName,
	}
	tasks, err := l.OtDB.FindByFilter(l.ctx, filter, nil)
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfo OtaTaskByDeviceName failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	var taskInfo dm.OtaUpTaskInfo
	_ = copier.Copy(&taskInfo, &tasks[0])
	return &taskInfo, nil
}
