package menu

import (
	"context"
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

func (l *IndexLogic) Index(req *types.MenuInfoIndexReq) (resp *types.MenuInfoIndexResp, err error) {
	info, err := l.svcCtx.ModuleRpc.ModuleMenuIndex(l.ctx, &sys.MenuInfoIndexReq{
		ModuleCode: req.ModuleCode,
		IsRetTree:  req.IsRetTree,
	})
	if err != nil {
		return nil, err
	}

	return &types.MenuInfoIndexResp{List: system.ToMenuInfosApi(info.List)}, nil
}
