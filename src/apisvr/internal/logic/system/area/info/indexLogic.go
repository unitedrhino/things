package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
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

func (l *IndexLogic) Index(req *types.AreaInfoIndexReq) (resp *types.AreaInfoIndexResp, err error) {
	dmReq := &sys.AreaInfoIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		ProjectID:  req.ProjectID,
		ProjectIDs: req.ProjectIDs,
		AreaID:     req.AreaID,
		AreaIDs:    req.AreaIDs,
	}
	dmResp, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AreaManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	list := make([]*types.AreaInfo, 0, len(dmResp.List))
	for _, pb := range dmResp.List {
		list = append(list, transPbToApi(pb))
	}

	return &types.AreaInfoIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
}
