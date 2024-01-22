package otaTask

import (
	"context"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobIndexLogic {
	return &JobIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobIndexLogic) JobIndex(req *types.OTATaskByJobIndexReq) (resp *types.OtaTaskByJobIndexResp, err error) {
	// todo: add your logic here and delete this line

	return
}
