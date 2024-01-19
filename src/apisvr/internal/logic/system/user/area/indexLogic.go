package area

import (
	"context"
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

func (l *IndexLogic) Index(req *types.UserAreaIndexReq) (resp *types.UserAreaIndexResp, err error) {
	dto := &sys.UserAreaIndexReq{
		Page:      logic.ToSysPageRpc(req.Page),
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
	}
	dmResp, err := l.svcCtx.UserRpc.UserAreaIndex(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.UserAreaIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	list := ToAreaApis(dmResp.List)
	return &types.UserAreaIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
}
