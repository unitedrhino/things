package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.ProjectInfoIndexReq) (resp *types.ProjectInfoIndexResp, err error) {
	dmReq := &sys.ProjectInfoIndexReq{
		Page:        logic.ToSysPageRpc(req.Page),
		ProjectName: req.ProjectName,
		ProjectIDs:  req.ProjectIDs,
	}
	dmResp, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProjectManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	list := make([]*types.ProjectInfo, 0, len(dmResp.List))
	for _, pb := range dmResp.List {
		list = append(list, system.ProjectInfoToApi(pb))
	}

	return &types.ProjectInfoIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
}
