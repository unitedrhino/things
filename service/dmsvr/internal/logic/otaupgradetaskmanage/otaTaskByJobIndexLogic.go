package otaupgradetaskmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskByJobIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OtDB *relationDB.OtaUpgradeTaskRepo
}

func NewOtaTaskByJobIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskByJobIndexLogic {
	return &OtaTaskByJobIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
	}
}

// 查询指定升级批次下的设备升级作业列表
func (l *OtaTaskByJobIndexLogic) OtaTaskByJobIndex(in *dm.OTATaskByJobIndexReq) (*dm.OtaTaskByJobIndexResp, error) {
	filter := relationDB.OtaUpgradeTaskFilter{JobId: in.JobId, DeviceNames: in.DeviceName}
	total, err := l.OtDB.CountByFilter(l.ctx, filter)
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfoIndexCount failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	var otaTaskInfo []*dm.OtaUpTaskInfo
	otaTask, err := l.OtDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.PageInfo))
	if err != nil {
		l.Errorf("%s.TaskInfo.TaskInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	_ = copier.Copy(&otaTaskInfo, &otaTask)
	return &dm.OtaTaskByJobIndexResp{Total: total, OtaUpTaskInfo: otaTaskInfo}, nil
}
