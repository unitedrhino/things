package otajobmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOTAStrategyByJobLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OjDB *relationDB.OtaJobRepo
}

func NewCancelOTAStrategyByJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOTAStrategyByJobLogic {
	return &CancelOTAStrategyByJobLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
	}
}

// 取消动态升级策略
func (l *CancelOTAStrategyByJobLogic) CancelOTAStrategyByJob(in *dm.JobReq) (*dm.Response, error) {
	otaJob, err := l.OjDB.FindOne(l.ctx, in.JobId)
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoRead failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	otaJob.JobStatus = msgOta.JobStatusCanceled
	err = l.OjDB.Update(l.ctx, otaJob)
	if err != nil {
		l.Errorf("%s.JobInfo.JobInfoUpdate failure err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
