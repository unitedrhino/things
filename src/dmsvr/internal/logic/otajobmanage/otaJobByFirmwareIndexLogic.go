package otajobmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaJobByFirmwareIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB *relationDB.OtaJobRepo
}

func NewOtaJobByFirmwareIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaJobByFirmwareIndexLogic {
	return &OtaJobByFirmwareIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// 获取升级包下的升级任务批次列表
func (l *OtaJobByFirmwareIndexLogic) OtaJobByFirmwareIndex(in *dm.OtaJobByFirmwareIndexReq) (*dm.OtaJobInfoIndexResp, error) {
	jobFilter := relationDB.OtaJobFilter{
		FirmwareId: in.FirmwareId,
	}
	total, err := l.OjDB.CountByFilter(l.ctx, jobFilter)
	if err != nil {
		return nil, err
	}
	var otaJobInfoList []*dm.OtaJobInfo
	otaJobList, err := l.OjDB.FindByFilter(l.ctx, jobFilter, logic.ToPageInfo(in.PageInfo))
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	err = copier.Copy(&otaJobInfoList, &otaJobList)
	if err != nil {
		l.Errorf("%s.JobInfo.CopyJobInfo failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.OtaJobInfoIndexResp{OtaJobInfo: otaJobInfoList, Total: total}, nil
}
