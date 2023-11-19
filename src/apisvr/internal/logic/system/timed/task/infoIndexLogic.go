package task

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoIndexLogic {
	return &InfoIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoIndexLogic) InfoIndex(req *types.TimedTaskInfoIndexReq) (resp *types.TimedTaskInfoIndexResp, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskInfoIndex(l.ctx, &timedjob.TaskInfoIndexReq{Page: logic.ToTimedJobPageRpc(req.Page)})
	if err != nil {
		return nil, err
	}
	return &types.TimedTaskInfoIndexResp{List: ToTaskInfosTypes(ret.List), Total: ret.Total}, nil
}
