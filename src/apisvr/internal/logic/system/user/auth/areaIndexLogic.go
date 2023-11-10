package auth

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaIndexLogic {
	return &AreaIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaIndexLogic) AreaIndex(req *types.UserAuthAreaIndexReq) (resp *types.UserAuthAreaIndexResp, err error) {
	dto := &sys.UserAuthAreaIndexReq{
		Page:      logic.ToSysPageRpc(req.Page),
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
	}
	dmResp, err := l.svcCtx.UserRpc.UserAuthAreaIndex(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.UserAuthAreaIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	list := ToAreaApis(dmResp.List)
	return &types.UserAuthAreaIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
	return
}
