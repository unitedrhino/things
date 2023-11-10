package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProjectInfoReadReq) (resp *types.ProjectInfo, err error) {
	dmResp, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectInfoReadReq{ProjectID: req.ProjectID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ProjectManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return system.ProjectInfoToApi(dmResp), nil
}
