package access

import (
	"context"
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

func (l *IndexLogic) Index(req *types.RoleAccessIndexReq) (resp *types.RoleAccessIndexResp, err error) {
	rt, err := l.svcCtx.RoleRpc.RoleAccessIndex(l.ctx, &sys.RoleAccessIndexReq{Id: req.ID})
	if err != nil {
		return nil, err
	}

	return &types.RoleAccessIndexResp{AccessCodes: rt.AccessCodes}, nil
}
