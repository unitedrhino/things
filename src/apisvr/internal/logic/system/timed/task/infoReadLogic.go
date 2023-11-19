package task

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoReadLogic {
	return &InfoReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoReadLogic) InfoRead(req *types.CodeReq) (resp *types.TimedTaskInfo, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskInfoRead(l.ctx, &timedjob.CodeReq{Code: req.Code})
	return ToTaskInfoTypes(ret), err
}
