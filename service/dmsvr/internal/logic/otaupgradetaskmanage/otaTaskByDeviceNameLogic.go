package otaupgradetaskmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskByDeviceNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaFirmwareDeviceRepo
}

func NewOtaTaskByDeviceNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskByDeviceNameLogic {
	return &OtaTaskByDeviceNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaFirmwareDeviceRepo(ctx),
	}
}

//
//// 根据设备的name和productId，查询task
//func (l *OtaTaskByDeviceNameLogic) OtaTaskByDeviceName(in *dm.OTATaskByDeviceNameReq) (*dm.OtaUpTaskInfo, error) {
//	filter := relationDB.OtaFirmwareDeviceFilter{
//		ProductID:  in.ProductID,
//		DeviceName: in.DeviceName,
//	}
//	tasks, err := l.OtDB.FindByFilter(l.ctx, filter, nil)
//	if err != nil {
//		l.Errorf("%s.TaskInfo.TaskInfo OtaTaskByDeviceName failure err=%+v", utils.FuncName(), err)
//		return nil, err
//	}
//	var taskInfo dm.OtaUpTaskInfo
//	_ = utils.CopyE(&taskInfo, &tasks[0])
//	return &taskInfo, nil
//}
