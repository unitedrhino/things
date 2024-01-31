package task

import (
	"context"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnalysisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAnalysisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalysisLogic {
	return &AnalysisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AnalysisLogic) Analysis(req *types.OtaTaskAnalysisReq) (resp *types.OtaTaskAnalysisResp, err error) {
	// todo: add your logic here and delete this line

	return
}
