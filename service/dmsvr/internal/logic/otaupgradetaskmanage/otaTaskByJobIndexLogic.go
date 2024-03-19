package otaupgradetaskmanagelogic

import (
	"context"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskByJobIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaFirmwareDeviceRepo
}

func NewOtaTaskByJobIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskByJobIndexLogic {
	return &OtaTaskByJobIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaFirmwareDeviceRepo(ctx),
	}
}

//
//// 查询指定升级批次下的设备升级作业列表
//func (l *OtaTaskByJobIndexLogic) OtaTaskByJobIndex(in *dm.OTATaskByJobIndexReq) (*dm.OtaTaskByJobIndexResp, error) {
//	filter := relationDB.OtaFirmwareDeviceFilter{JobID: in.JobID, DeviceNames: in.DeviceNames}
//	total, err := l.OtDB.CountByFilter(l.ctx, filter)
//	if err != nil {
//		l.Errorf("%s.TaskInfo.TaskInfoIndexCount failure err=%+v", utils.FuncName(), err)
//		return nil, err
//	}
//	var otaTaskInfo []*dm.OtaUpTaskInfo
//	otaTask, err := l.OtDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.PageInfo))
//	if err != nil {
//		l.Errorf("%s.TaskInfo.TaskInfo failure err=%+v", utils.FuncName(), err)
//		return nil, err
//	}
//	_ = copier.Copy(&otaTaskInfo, &otaTask)
//	return &dm.OtaTaskByJobIndexResp{Total: total, List: otaTaskInfo}, nil
//}
