package otaTask

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobCancelLogic {
	return &JobCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobCancelLogic) JobCancel(req *types.OTATaskByJobCancelReq) error {
	// todo: add your logic here and delete this line

	return nil
}
